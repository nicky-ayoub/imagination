package viewport

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

//-------------------------------------------------------------------------------------------------
// Types
//-------------------------------------------------------------------------------------------------
/** All available flipping modes. */
type TViewportFlippingModeID int

const (
	FLIPPING_MODE_ID_NORMAL                  TViewportFlippingModeID = iota //!< Flipping is disabled.
	FLIPPING_MODE_ID_HORIZONTAL                                             //!< Horizontal flipping.
	FLIPPING_MODE_ID_VERTICAL                                               //!< Vertical flipping.
	FLIPPING_MODE_ID_HORIZONTAL_AND_VERTICAL                                //!< Both horizontal and vertical flipping.
	FLIPPING_MODE_IDS_COUNT                                                 //!< How many flipping modes are available.
)

//-------------------------------------------------------------------------------------------------
// Constants
//-------------------------------------------------------------------------------------------------
const (
	/** The display refresh rate period in milliseconds (1 / 60Hz ~= 16ms). 60 frames per second are fine for a video game, so they are more than enough for an image viewer. */
	DISPLAY_REFRESH_RATE_PERIOD uint32 = 16
	/** The maximum allowed zoom factor value. */
	VIEWPORT_MAXIMUM_ZOOM_FACTOR int32 = 16
	/** Window minimum width in pixels. */
	VIEWPORT_MINIMUM_WINDOW_WIDTH int32 = 100
	/** Window minimum height in pixels. */
	VIEWPORT_MINIMUM_WINDOW_HEIGHT int32 = 100
)

type viewport struct {
	width                 int32
	height                int32
	original_width        int32
	original_height       int32
	original_pixel_format uint32
	adjusted_width        int32
	adjusted_height       int32
	flip_mode             sdl.RendererFlip
	srcRect               sdl.Rect
}

var vp viewport

var window *sdl.Window
var renderer *sdl.Renderer

/** The texture holding the loaded image. */
var texture *sdl.Texture = nil

/** The texture holding the image adapted to the current viewport dimensions. */
var adapted_texture *sdl.Texture = nil

func init() {
	vp.flip_mode = sdl.FLIP_NONE
}

//-------------------------------------------------------------------------------------------------
// Private functions
//-------------------------------------------------------------------------------------------------
/** Add eventual additional borders to the original image to make sure its ratio is kept regardless of the viewport dimensions.
 * @param Image_Width The image to display width in pixels.
 * @param Image_Height The image to display height in pixels.
 * @return 0 if the function succeeded,
 * @return -1 if an error occurred.
 */
func adaptImage(width int32, height int32) (err error) {
	var Horizontal_Scaling_Percentage int32
	var Vertical_Scaling_Percentage int32
	var dstRect sdl.Rect

	//fmt.Println("AdaptImage(", Image_Width, ", ", Image_Height, ")")
	//fmt.Printf("Viewport(%d, %d)\n", vp.width, vp.height)

	// Adjust image size to viewport ratio to make sure the image will keep its ratio
	// Image is smaller than the viewport, keep viewport size
	if (width < vp.width) && (height < vp.height) {
		vp.adjusted_width = vp.width
		vp.adjusted_height = vp.height
	} else { // Image is bigger than the viewport, make sure it will fit in the current viewport
		// Determine which size (horizontal or vertical) must be filled to keep the ratio (all computations are made with percentages to be easier to understand, and also because multiplying by 100 allows simple and fast fixed calculations with accurate enough results)
		Horizontal_Scaling_Percentage = (100 * width) / vp.width // Find how much the image is larger than the viewport
		Vertical_Scaling_Percentage = (100 * height) / vp.height // Find how much the image is higher than the viewport

		// Horizontal size is more "compressed" than vertical one for this viewport, so display the full image horizontal size and adjust vertical size to keep ratio
		if Horizontal_Scaling_Percentage > Vertical_Scaling_Percentage {
			vp.adjusted_width = width                                              // Keep the width
			vp.adjusted_height = (vp.height * Horizontal_Scaling_Percentage) / 100 // Fill the image bottom to keep the ratio
		} else { // Vertical size is more "compressed" than horizontal one for this viewport, so display the full image vertical size and adjust horizontal size to keep ratio
			vp.adjusted_height = height                                        // Keep the height
			vp.adjusted_width = (vp.width * Vertical_Scaling_Percentage) / 100 // Fill the image right to keep the ratio
		}
	}

	// Remove the previously existing texture
	if adapted_texture != nil {
		adapted_texture.Destroy()
	}

	// Create a texture with the right dimensions to keep the image ratio
	adapted_texture, err = renderer.CreateTexture(vp.original_pixel_format, sdl.TEXTUREACCESS_TARGET, vp.adjusted_width, vp.adjusted_height)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Fill the texture with a visible background color, so the original image dimensions are easily visible

	if err = renderer.SetRenderTarget(adapted_texture); err != nil {
		log.Fatal(err)
		return err
	}
	// Set a black background color
	if err = renderer.SetDrawColor(192, 192, 192, 0); err != nil {
		log.Fatal(err)
		return err
	}
	// Do the fill operation
	if err = renderer.Clear(); err != nil {
		log.Fatal(err)
		return err
	}

	// Copy the original image on the adjusted image
	dstRect.X = renderer.GetViewport().W/2 - width/2
	dstRect.Y = renderer.GetViewport().H/2 - height/2
	dstRect.W = width - 1
	dstRect.H = height - 1

	if err = renderer.CopyEx(texture, nil, &dstRect, 0, nil, vp.flip_mode); err != nil {
		log.Fatal(err)
		return err
	}

	// Restore the renderer (it can write to the display again)
	if err = renderer.SetRenderTarget(nil); err != nil {
		log.Fatal(err)
		return err
	}

	// Reset zoom when resizing the window to avoid aiming to whatever but the right place
	SetZoomedArea(0, 0, 1)

	return
}

//-------------------------------------------------------------------------------------------------
// Public functions
//-------------------------------------------------------------------------------------------------
func Initialize(title string, image *sdl.Surface) (err error) {
	// Try to create the viewport window
	if window == nil {
		if window, err = sdl.CreateWindow("-", 0, 0, 640, 480, sdl.WINDOW_RESIZABLE|sdl.WINDOW_MAXIMIZED); err != nil {
			log.Fatal(err)
			return err
		}
	}
	window.SetTitle(title)

	// Try to create an hardware-accelerated renderer to plug to the window
	if renderer == nil {
		if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC); err != nil {
			log.Fatal(err)
			return err
		}
	}
	if texture != nil {
		texture.Destroy()
	}
	// Convert the image surface to a texture
	if texture, err = renderer.CreateTextureFromSurface(image); err != nil {
		log.Fatal(err)
		return err
	}

	// Do not allow the window to be too small because it can prevent the texture rendering from working
	window.SetMinimumSize(VIEWPORT_MINIMUM_WINDOW_WIDTH, VIEWPORT_MINIMUM_WINDOW_HEIGHT)

	// Cache original image dimensions
	if vp.original_pixel_format, _, vp.original_width, vp.original_height, err = texture.Query(); err != nil {
		log.Fatal(err)
		return err
	}
	//fmt.Printf("Initialize():  Image W=%d, H=%d\n", vp.original_image_width, vp.original_image_height)

	return nil
}

func DrawImage() {
	renderer.CopyEx(adapted_texture, &vp.srcRect, nil, 0, nil, 0)
	renderer.Present()
}

func SetDimensions(new_width int32, new_height int32) {
	// Store new viewport dimensions
	vp.width = new_width
	vp.height = new_height

	// Add additional borders to the image to keep its ratio
	adaptImage(vp.original_width, vp.original_height)
}

var Previous_Zoom_Level_Rectangle_X = int32(0)
var Previous_Zoom_Level_Rectangle_Y = int32(0)
var Previous_Zoom_Factor = int32(1)

func SetZoomedArea(Viewport_X int32, Viewport_Y int32, Zoom_Factor int32) {
	var Rectangle_X = int32(0)
	var Rectangle_Y = int32(0)

	// Do not compute viewing area once more if the maximum zooming level has been reached, because values would overflow
	if (Previous_Zoom_Factor == VIEWPORT_MAXIMUM_ZOOM_FACTOR) && (Zoom_Factor == VIEWPORT_MAXIMUM_ZOOM_FACTOR) {
		return
	}

	// Force the view rectangle to start from the viewport origin when there is no zooming
	if Zoom_Factor == 1 {
		Rectangle_X = 0
		Rectangle_Y = 0
	} else if Previous_Zoom_Factor < Zoom_Factor {
		// Handle zooming in by adding to the preceding view rectangle origin the new mouse moves (scaled according to the new zoom factor)
		Rectangle_X = Previous_Zoom_Level_Rectangle_X + (((Viewport_X / Zoom_Factor) * vp.adjusted_width) / vp.width)
		Rectangle_Y = Previous_Zoom_Level_Rectangle_Y + (((Viewport_Y / Zoom_Factor) * vp.adjusted_height) / vp.height)
	} else if Previous_Zoom_Factor > Zoom_Factor {
		// Handle zooming out by subtracting to the preceding view rectangle origin the new mouse moves (scaled according to the previous zoom factor, which was greater than the current one and was the factor used to compute the zooming in)
		Rectangle_X = Previous_Zoom_Level_Rectangle_X - (((Viewport_X / Previous_Zoom_Factor) * vp.adjusted_width) / vp.width)
		Rectangle_Y = Previous_Zoom_Level_Rectangle_Y - (((Viewport_Y / Previous_Zoom_Factor) * vp.adjusted_height) / vp.height)
	}

	// Make sure no negative coordinates are generated
	if Rectangle_X < 0 {
		Rectangle_X = 0
	}
	if Rectangle_Y < 0 {
		Rectangle_Y = 0
	}

	// There is nothing to do when zooming to 1x because the for loop will immediately exit and x and y coordinates will be set to 0
	vp.srcRect.X = Rectangle_X
	vp.srcRect.Y = Rectangle_Y

	// The smaller the rectangle is, the more the image will be zoomed
	vp.srcRect.W = (vp.adjusted_width / Zoom_Factor) - 1
	vp.srcRect.H = (vp.adjusted_height / Zoom_Factor) - 1

	Previous_Zoom_Level_Rectangle_X = Rectangle_X
	Previous_Zoom_Level_Rectangle_Y = Rectangle_Y
	Previous_Zoom_Factor = Zoom_Factor
}

func SetFlippingMode(mode TViewportFlippingModeID) {
	// Update renderer flip flags according to the new mode
	switch mode {
	case FLIPPING_MODE_ID_NORMAL:
		vp.flip_mode = sdl.FLIP_NONE
	case FLIPPING_MODE_ID_HORIZONTAL:
		vp.flip_mode = sdl.FLIP_HORIZONTAL
	case FLIPPING_MODE_ID_VERTICAL:
		vp.flip_mode = sdl.FLIP_VERTICAL
	case FLIPPING_MODE_ID_HORIZONTAL_AND_VERTICAL:
		vp.flip_mode = sdl.FLIP_HORIZONTAL | sdl.FLIP_VERTICAL
	default:
		log.Fatal("Error : bad flipping mode ID provided")
		return
	}

	// Redraw the image with the newly selected flipping mode
	adaptImage(vp.original_width, vp.original_height)
}

func ScaleImage() {
	var Horizontally_Scaled_Pixels_Count int32
	var Vertically_Scaled_Pixels_Count int32
	var Scaling_Percentage int32

	// Always reset the zoom to ease the following computations
	SetZoomedArea(0, 0, 1)

	// Make the image fit the viewport if it is smaller
	if (vp.original_width < vp.width) && (vp.original_height < vp.height) {
		// Determine the amount of pixels not used by the image and keep the smallest one to make sure the image ratio is not modified
		Horizontally_Scaled_Pixels_Count = vp.width - vp.original_width
		Vertically_Scaled_Pixels_Count = vp.height - vp.original_height
		if Horizontally_Scaled_Pixels_Count < Vertically_Scaled_Pixels_Count {
			// Compute how many percents the image will be horizontally scaled
			Scaling_Percentage = (100 * (vp.original_width + Horizontally_Scaled_Pixels_Count)) / vp.original_width

			// Scale the "camera"
			vp.srcRect.W = vp.original_width + Horizontally_Scaled_Pixels_Count
			vp.srcRect.H = (vp.original_height * Scaling_Percentage) / 100 // Use the percentage computed right before to scale the vertical direction with the same proportion
		} else {
			// Compute how many percents the image will be horizontally scaled
			Scaling_Percentage = (100 * (vp.original_height + Vertically_Scaled_Pixels_Count)) / vp.original_height

			// Scale the "camera"
			vp.srcRect.W = (vp.original_width * Scaling_Percentage) / 100
			vp.srcRect.H = vp.original_height + Vertically_Scaled_Pixels_Count
		}

		// Make sure the camera will display the whole image
		vp.srcRect.X = 0
		vp.srcRect.Y = 0

		// Fill the empty part of the image if its ratio is different from the viewport ratio
		adaptImage(vp.srcRect.W, vp.srcRect.H)
	}
}
