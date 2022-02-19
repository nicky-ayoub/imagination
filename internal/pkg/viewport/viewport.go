package viewport

import (
	"fmt"
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
	VIEWPORT_MAXIMUM_ZOOM_FACTOR int32 = 256
	/** Window minimum width in pixels. */
	VIEWPORT_MINIMUM_WINDOW_WIDTH int32 = 100
	/** Window minimum height in pixels. */
	VIEWPORT_MINIMUM_WINDOW_HEIGHT int32 = 100
)

var window *sdl.Window

/** The renderer able to display to the window. */
var renderer *sdl.Renderer

/** Viewport width in pixels. */
var Viewport_Width int32

/** Viewport height in pixels. */
var Viewport_Height int32

/** Cache the loaded image width in pixels. */
var Viewport_Original_Image_Width int32

/** Cache the loaded image height in pixels. */
var Viewport_Original_Image_Height int32

/** Cache the loaded image pixel format. */
var Viewport_Original_Image_Pixel_Format uint32

/** Adjusted image width in pixels. */
var Viewport_Adjusted_Image_Width int32

/** Adjusted image height in pixels. */
var Viewport_Adjusted_Image_Height int32

/** The texture holding the loaded image. */
var texture *sdl.Texture = nil

/** The texture holding the image adapted to the current viewport dimensions. */
var adapted_texture *sdl.Texture = nil

/** The rectangle defining the area of the adjusted image to display to the viewport. */
var Viewport_Rectangle_View sdl.Rect

/** The flip value to apply to the image when adapting it to the viewport. */
var Viewport_Adapted_Image_Flipping_Mode sdl.RendererFlip = sdl.FLIP_NONE

//-------------------------------------------------------------------------------------------------
// Private functions
//-------------------------------------------------------------------------------------------------
/** Add eventual additional borders to the original image to make sure its ratio is kept regardless of the viewport dimensions.
 * @param Image_Width The image to display width in pixels.
 * @param Image_Height The image to display height in pixels.
 * @return 0 if the function succeeded,
 * @return -1 if an error occurred.
 */
func AdaptImage(Image_Width int32, Image_Height int32) (err error) {
	var Horizontal_Scaling_Percentage int32
	var Vertical_Scaling_Percentage int32
	var Rectangle_Original_Image_Dimensions sdl.Rect

	//fmt.Println("AdaptImage(", Image_Width, ", ", Image_Height, ")")
	//fmt.Printf("Viewport(%d, %d)\n", Viewport_Width, Viewport_Height)

	// Adjust image size to viewport ratio to make sure the image will keep its ratio
	// Image is smaller than the viewport, keep viewport size
	if (Image_Width < Viewport_Width) && (Image_Height < Viewport_Height) {
		Viewport_Adjusted_Image_Width = Viewport_Width
		Viewport_Adjusted_Image_Height = Viewport_Height
	} else { // Image is bigger than the viewport, make sure it will fit in the current viewport
		// Determine which size (horizontal or vertical) must be filled to keep the ratio (all computations are made with percentages to be easier to understand, and also because multiplying by 100 allows simple and fast fixed calculations with accurate enough results)
		Horizontal_Scaling_Percentage = (100 * Image_Width) / Viewport_Width // Find how much the image is larger than the viewport
		Vertical_Scaling_Percentage = (100 * Image_Height) / Viewport_Height // Find how much the image is higher than the viewport

		// Horizontal size is more "compressed" than vertical one for this viewport, so display the full image horizontal size and adjust vertical size to keep ratio
		if Horizontal_Scaling_Percentage > Vertical_Scaling_Percentage {
			Viewport_Adjusted_Image_Width = Image_Width                                              // Keep the width
			Viewport_Adjusted_Image_Height = (Viewport_Height * Horizontal_Scaling_Percentage) / 100 // Fill the image bottom to keep the ratio
		} else { // Vertical size is more "compressed" than horizontal one for this viewport, so display the full image vertical size and adjust horizontal size to keep ratio
			Viewport_Adjusted_Image_Height = Image_Height                                        // Keep the height
			Viewport_Adjusted_Image_Width = (Viewport_Width * Vertical_Scaling_Percentage) / 100 // Fill the image right to keep the ratio
		}
	}

	// Remove the previously existing texture
	if adapted_texture != nil {
		adapted_texture.Destroy()
	}

	// Create a texture with the right dimensions to keep the image ratio
	adapted_texture, err = renderer.CreateTexture(Viewport_Original_Image_Pixel_Format, sdl.TEXTUREACCESS_TARGET, Viewport_Adjusted_Image_Width, Viewport_Adjusted_Image_Height)
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
	if err = renderer.SetDrawColor(0, 0, 0, 0); err != nil {
		log.Fatal(err)
		return err
	}
	// Do the fill operation
	if err = renderer.Clear(); err != nil {
		log.Fatal(err)
		return err
	}

	// Copy the original image on the adjusted image
	Rectangle_Original_Image_Dimensions.X = renderer.GetViewport().W/2 - Image_Width/2
	Rectangle_Original_Image_Dimensions.Y = renderer.GetViewport().H/2 - Image_Height/2
	Rectangle_Original_Image_Dimensions.W = Image_Width - 1
	Rectangle_Original_Image_Dimensions.H = Image_Height - 1

	err = renderer.CopyEx(texture, nil, &Rectangle_Original_Image_Dimensions, 0, nil, Viewport_Adapted_Image_Flipping_Mode)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Restore the renderer (it can write to the display again)
	err = renderer.SetRenderTarget(nil)
	if err != nil {
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
	window, err = sdl.CreateWindow(title, 0, 0, 640, 480, sdl.WINDOW_RESIZABLE|sdl.WINDOW_MAXIMIZED)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Try to create an hardware-accelerated renderer to plug to the window
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Convert the image surface to a texture
	texture, err = renderer.CreateTextureFromSurface(image)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Do not allow the window to be too small because it can prevent the texture rendering from working
	window.SetMinimumSize(VIEWPORT_MINIMUM_WINDOW_WIDTH, VIEWPORT_MINIMUM_WINDOW_HEIGHT)

	// Cache original image dimensions
	Viewport_Original_Image_Pixel_Format, _, Viewport_Original_Image_Width, Viewport_Original_Image_Height, err = texture.Query()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("Initialize():  Image W=%d, H=%d\n", Viewport_Original_Image_Width, Viewport_Original_Image_Height)

	return nil
}

func DrawImage() {
	renderer.CopyEx(adapted_texture, &Viewport_Rectangle_View, nil, 0, nil, 0)
	renderer.Present()
}

func SetDimensions(New_Viewport_Width int32, New_Viewport_Height int32) {
	// Store new viewport dimensions
	Viewport_Width = New_Viewport_Width
	Viewport_Height = New_Viewport_Height

	// Add additional borders to the image to keep its ratio
	AdaptImage(Viewport_Original_Image_Width, Viewport_Original_Image_Height)
}

var SetZoomedArea func(int32, int32, int32)

func init() {
	//fmt.Println("Creating SetZoomedArea function...")
	SetZoomedArea = MakeSetZoomedArea()
}

func MakeSetZoomedArea() func(int32, int32, int32) {
	var Previous_Zoom_Level_Rectangle_X = int32(0)
	var Previous_Zoom_Level_Rectangle_Y = int32(0)
	var Previous_Zoom_Factor = int32(1)

	return func(Viewport_X int32, Viewport_Y int32, Zoom_Factor int32) {
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
			Rectangle_X = Previous_Zoom_Level_Rectangle_X + (((Viewport_X / Zoom_Factor) * Viewport_Adjusted_Image_Width) / Viewport_Width)
			Rectangle_Y = Previous_Zoom_Level_Rectangle_Y + (((Viewport_Y / Zoom_Factor) * Viewport_Adjusted_Image_Height) / Viewport_Height)
		} else if Previous_Zoom_Factor > Zoom_Factor {
			// Handle zooming out by subtracting to the preceding view rectangle origin the new mouse moves (scaled according to the previous zoom factor, which was greater than the current one and was the factor used to compute the zooming in)
			Rectangle_X = Previous_Zoom_Level_Rectangle_X - (((Viewport_X / Previous_Zoom_Factor) * Viewport_Adjusted_Image_Width) / Viewport_Width)
			Rectangle_Y = Previous_Zoom_Level_Rectangle_Y - (((Viewport_Y / Previous_Zoom_Factor) * Viewport_Adjusted_Image_Height) / Viewport_Height)
		}

		// Make sure no negative coordinates are generated
		if Rectangle_X < 0 {
			Rectangle_X = 0
		}
		if Rectangle_Y < 0 {
			Rectangle_Y = 0
		}

		// There is nothing to do when zooming to 1x because the for loop will immediately exit and x and y coordinates will be set to 0
		Viewport_Rectangle_View.X = Rectangle_X
		Viewport_Rectangle_View.Y = Rectangle_Y

		// The smaller the rectangle is, the more the image will be zoomed
		Viewport_Rectangle_View.W = (Viewport_Adjusted_Image_Width / Zoom_Factor) - 1
		Viewport_Rectangle_View.H = (Viewport_Adjusted_Image_Height / Zoom_Factor) - 1

		Previous_Zoom_Level_Rectangle_X = Rectangle_X
		Previous_Zoom_Level_Rectangle_Y = Rectangle_Y
		Previous_Zoom_Factor = Zoom_Factor
	}
}

func SetFlippingMode(Flipping_Mode_ID TViewportFlippingModeID) {
	// Update renderer flip flags according to the new mode
	switch Flipping_Mode_ID {
	case FLIPPING_MODE_ID_NORMAL:
		Viewport_Adapted_Image_Flipping_Mode = sdl.FLIP_NONE
	case FLIPPING_MODE_ID_HORIZONTAL:
		Viewport_Adapted_Image_Flipping_Mode = sdl.FLIP_HORIZONTAL
	case FLIPPING_MODE_ID_VERTICAL:
		Viewport_Adapted_Image_Flipping_Mode = sdl.FLIP_VERTICAL
	case FLIPPING_MODE_ID_HORIZONTAL_AND_VERTICAL:
		Viewport_Adapted_Image_Flipping_Mode = sdl.FLIP_HORIZONTAL | sdl.FLIP_VERTICAL
	default:
		log.Fatal("Error : bad flipping mode ID provided")
		return
	}

	// Redraw the image with the newly selected flipping mode
	AdaptImage(Viewport_Original_Image_Width, Viewport_Original_Image_Height)
}

func ScaleImage() {
	var Horizontally_Scaled_Pixels_Count int32
	var Vertically_Scaled_Pixels_Count int32
	var Scaling_Percentage int32

	// Always reset the zoom to ease the following computations
	SetZoomedArea(0, 0, 1)

	// Make the image fit the viewport if it is smaller
	if (Viewport_Original_Image_Width < Viewport_Width) && (Viewport_Original_Image_Height < Viewport_Height) {
		// Determine the amount of pixels not used by the image and keep the smallest one to make sure the image ratio is not modified
		Horizontally_Scaled_Pixels_Count = Viewport_Width - Viewport_Original_Image_Width
		Vertically_Scaled_Pixels_Count = Viewport_Height - Viewport_Original_Image_Height
		if Horizontally_Scaled_Pixels_Count < Vertically_Scaled_Pixels_Count {
			// Compute how many percents the image will be horizontally scaled
			Scaling_Percentage = (100 * (Viewport_Original_Image_Width + Horizontally_Scaled_Pixels_Count)) / Viewport_Original_Image_Width

			// Scale the "camera"
			Viewport_Rectangle_View.W = Viewport_Original_Image_Width + Horizontally_Scaled_Pixels_Count
			Viewport_Rectangle_View.H = (Viewport_Original_Image_Height * Scaling_Percentage) / 100 // Use the percentage computed right before to scale the vertical direction with the same proportion
		} else {
			// Compute how many percents the image will be horizontally scaled
			Scaling_Percentage = (100 * (Viewport_Original_Image_Height + Vertically_Scaled_Pixels_Count)) / Viewport_Original_Image_Height

			// Scale the "camera"
			Viewport_Rectangle_View.W = (Viewport_Original_Image_Width * Scaling_Percentage) / 100
			Viewport_Rectangle_View.H = Viewport_Original_Image_Height + Vertically_Scaled_Pixels_Count
		}

		// Make sure the camera will display the whole image
		Viewport_Rectangle_View.X = 0
		Viewport_Rectangle_View.Y = 0

		// Fill the empty part of the image if its ratio is different from the viewport ratio
		AdaptImage(Viewport_Rectangle_View.W, Viewport_Rectangle_View.H)
	}
}
