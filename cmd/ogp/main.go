// Command ogp generates the static background used for per-article OGP
// (Twitter Card) images. It renders a themed gradient, a translucent
// "glass" panel, an accent bar, and the bottom-left branding (avatar +
// site name). The article title itself is overlaid later by Hugo at build
// time via the images.Text filter, so this background is fixed and shared
// by every post.
//
// Usage:
//
//	go run ./cmd/ogp
//
// It writes assets/ogp/base.png, which is committed so that CI does not
// need the Go toolchain to build OGP images.
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	_ "golang.org/x/image/webp"
)

const (
	width  = 1200
	height = 630

	margin       = 72  // outer margin of the glass panel
	cornerRadius = 36  // glass panel corner radius
	pad          = 48  // inner padding of the glass panel
	avatarSize   = 104 // diameter of the circular avatar

	fontPath   = "assets/ogp/NotoSansJP-Bold.ttf"
	avatarPath = "assets/image/avatar.webp"
	outPath    = "assets/ogp/base.png"
	siteName   = "りんりんのメモ帳"
)

// Theme palette mirrors the liquid-glass theme's dark scheme
// (assets/css/partials/_tokens.css and _scene-blob.css): a near-black
// navy base with soft, drifting accent "blobs" and translucent glass.
var (
	baseBG    = color.RGBA{0x0b, 0x0e, 0x1a, 0xff} // --color-bg (dark)
	glassFill = color.RGBA{0xff, 0xff, 0xff, 0x1f} // --glass-white  ~0.12
	glassEdge = color.RGBA{0xff, 0xff, 0xff, 0x40} // --glass-border ~0.25
	accentA   = color.RGBA{0x5e, 0xe7, 0xdf, 0xff} // --accent-aqua
	accentB   = color.RGBA{0xb4, 0x90, 0xf5, 0xff} // --accent-violet
	brandText = color.RGBA{0xff, 0xff, 0xff, 0xff} // --color-text (dark)
)

// blob describes one radial-gradient glow from the theme's background scene.
type blob struct {
	cx, cy, r float64
	c1, c2    color.RGBA // center -> edge, matching the theme's blobs
	opacity   float64
}

var blobs = []blob{
	// aqua -> blue, top-right
	{1060, 110, 540, color.RGBA{0x5e, 0xe7, 0xdf, 0xff}, color.RGBA{0x3b, 0x82, 0xf6, 0xff}, 0.55},
	// violet -> pink, bottom-left
	{170, 560, 580, color.RGBA{0xb4, 0x90, 0xf5, 0xff}, color.RGBA{0xec, 0x48, 0x99, 0xff}, 0.5},
	// amber -> rose, bottom-right
	{990, 640, 460, color.RGBA{0xff, 0xd2, 0x7f, 0xff}, color.RGBA{0xf7, 0xa8, 0xc4, 0xff}, 0.42},
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("ogp: %v", err)
	}
}

func run() error {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	drawScene(dst)

	panel := image.Rect(margin, margin, width-margin, height-margin)
	drawGlassPanel(dst, panel)

	// Accent bar (aqua -> violet) at the top-left of the panel.
	bar := image.Rect(panel.Min.X+pad, panel.Min.Y+pad, panel.Min.X+pad+72, panel.Min.Y+pad+10)
	fillBarGradient(dst, bar, 5, accentA, accentB)

	face, err := loadFace(30)
	if err != nil {
		return err
	}
	defer face.Close()

	// Branding: circular avatar + site name in the bottom-left of the panel.
	ax := panel.Min.X + pad
	ay := panel.Max.Y - pad - avatarSize
	if err := drawAvatar(dst, ax, ay); err != nil {
		return err
	}
	drawText(dst, face, siteName, ax+avatarSize+28, ay+avatarSize/2)

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, dst)
}

// drawScene paints the dark base then screen-blends the theme's soft accent
// blobs over it, reproducing the liquid-glass background scene. Each blob is
// only evaluated within its bounding box, and the (cheaper) squared distance
// is used to reject pixels outside its radius before taking a square root.
func drawScene(dst *image.RGBA) {
	bounds := dst.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.SetRGBA(x, y, baseBG)
		}
	}
	for _, b := range blobs {
		minX := max(bounds.Min.X, int(math.Floor(b.cx-b.r)))
		maxX := min(bounds.Max.X, int(math.Ceil(b.cx+b.r)))
		minY := max(bounds.Min.Y, int(math.Floor(b.cy-b.r)))
		maxY := min(bounds.Max.Y, int(math.Ceil(b.cy+b.r)))
		r2 := b.r * b.r
		for y := minY; y < maxY; y++ {
			for x := minX; x < maxX; x++ {
				dx := float64(x) + 0.5 - b.cx
				dy := float64(y) + 0.5 - b.cy
				d2 := dx*dx + dy*dy
				if d2 >= r2 {
					continue
				}
				t := math.Sqrt(d2) / b.r
				col := lerp(b.c1, b.c2, t)
				// Strong at the center, smoothly fading to nothing at the edge.
				a := b.opacity * smoothstep(1, 0, t)
				if a <= 0 {
					continue
				}
				dst.SetRGBA(x, y, screen(dst.RGBAAt(x, y), col, a))
			}
		}
	}
}

// drawGlassPanel overlays a translucent rounded rectangle with a soft border.
func drawGlassPanel(dst *image.RGBA, r image.Rectangle) {
	fillRoundRect(dst, r, cornerRadius, glassFill)
	strokeRoundRect(dst, r, cornerRadius, glassEdge, 2)
}

// drawAvatar decodes the site avatar, scales it, applies a circular mask
// with an anti-aliased edge, and composites it onto dst at (x, y).
func drawAvatar(dst *image.RGBA, x, y int) error {
	f, err := os.Open(avatarPath)
	if err != nil {
		return err
	}
	defer f.Close()
	src, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	scaled := image.NewRGBA(image.Rect(0, 0, avatarSize, avatarSize))
	draw.CatmullRom.Scale(scaled, scaled.Bounds(), src, src.Bounds(), draw.Over, nil)

	c := float64(avatarSize) / 2
	for j := 0; j < avatarSize; j++ {
		for i := 0; i < avatarSize; i++ {
			d := math.Hypot(float64(i)+0.5-c, float64(j)+0.5-c)
			a := c - d // distance inside the circle edge, for 1px AA
			if a <= 0 {
				continue
			}
			blendPixel(dst, x+i, y+j, scaled.RGBAAt(i, j), math.Min(a, 1))
		}
	}
	return nil
}

// drawText draws s with its vertical center at cy, left-aligned at x.
func drawText(dst *image.RGBA, face font.Face, s string, x, cy int) {
	m := face.Metrics()
	baseline := cy + (m.Ascent-m.Descent).Round()/2
	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(brandText),
		Face: face,
		Dot:  fixed.P(x, baseline),
	}
	d.DrawString(s)
}

func loadFace(size float64) (font.Face, error) {
	b, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	ft, err := opentype.Parse(b)
	if err != nil {
		return nil, err
	}
	return opentype.NewFace(ft, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
}

// --- low-level drawing helpers ---

func lerp(a, b color.RGBA, t float64) color.RGBA {
	t = clamp(t, 0, 1)
	return color.RGBA{
		R: uint8(float64(a.R) + (float64(b.R)-float64(a.R))*t),
		G: uint8(float64(a.G) + (float64(b.G)-float64(a.G))*t),
		B: uint8(float64(a.B) + (float64(b.B)-float64(a.B))*t),
		A: 0xff,
	}
}

// smoothstep returns a smooth 0..1 interpolation of x between edges e0 and e1.
func smoothstep(e0, e1, x float64) float64 {
	t := clamp((x-e0)/(e1-e0), 0, 1)
	return t * t * (3 - 2*t)
}

// screen blends light color c (scaled by alpha a) over bg using the screen
// operator, which brightens like overlapping glows on a dark background.
func screen(bg, c color.RGBA, a float64) color.RGBA {
	sc := func(b, s uint8) uint8 {
		bf, sf := float64(b)/255, float64(s)/255*a
		return uint8((1 - (1-bf)*(1-sf)) * 255)
	}
	return color.RGBA{sc(bg.R, c.R), sc(bg.G, c.G), sc(bg.B, c.B), 0xff}
}

// blendPixel alpha-composites c over dst at (x,y), scaling c's alpha by k (0..1).
func blendPixel(dst *image.RGBA, x, y int, c color.RGBA, k float64) {
	if !image.Pt(x, y).In(dst.Rect) {
		return
	}
	a := float64(c.A) / 255 * k
	if a <= 0 {
		return
	}
	bg := dst.RGBAAt(x, y)
	dst.SetRGBA(x, y, color.RGBA{
		R: uint8(float64(c.R)*a + float64(bg.R)*(1-a)),
		G: uint8(float64(c.G)*a + float64(bg.G)*(1-a)),
		B: uint8(float64(c.B)*a + float64(bg.B)*(1-a)),
		A: 0xff,
	})
}

// inRoundRect reports the coverage (0..1, anti-aliased) of point (px,py)
// inside rectangle r with corner radius rad.
func inRoundRect(px, py float64, r image.Rectangle, rad float64) float64 {
	minX, minY := float64(r.Min.X), float64(r.Min.Y)
	maxX, maxY := float64(r.Max.X), float64(r.Max.Y)
	// Nearest point on the inner rounded box.
	cx := clamp(px, minX+rad, maxX-rad)
	cy := clamp(py, minY+rad, maxY-rad)
	dx, dy := px-cx, py-cy
	// In the straight (non-corner) zones one of dx/dy is zero, so distance
	// reduces to edge distance and AA still works.
	if (px >= minX+rad && px <= maxX-rad) || (py >= minY+rad && py <= maxY-rad) {
		// Simple inside test for edges, with 1px AA on the outer border.
		edge := math.Min(math.Min(px-minX, maxX-px), math.Min(py-minY, maxY-py))
		return clamp(edge, 0, 1)
	}
	d := math.Hypot(dx, dy)
	return clamp(rad-d, 0, 1)
}

// fillRoundRectFunc fills a rounded rectangle, taking each pixel's color from
// colorAt (called with the column x), so callers can supply a solid color or a
// horizontal gradient without duplicating the coverage loop.
func fillRoundRectFunc(dst *image.RGBA, r image.Rectangle, rad int, colorAt func(x int) color.RGBA) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			cov := inRoundRect(float64(x)+0.5, float64(y)+0.5, r, float64(rad))
			if cov > 0 {
				blendPixel(dst, x, y, colorAt(x), cov)
			}
		}
	}
}

// fillRoundRect fills a rounded rectangle with a solid color.
func fillRoundRect(dst *image.RGBA, r image.Rectangle, rad int, c color.RGBA) {
	fillRoundRectFunc(dst, r, rad, func(int) color.RGBA { return c })
}

// fillBarGradient fills a small rounded bar with a horizontal c1->c2 gradient.
func fillBarGradient(dst *image.RGBA, r image.Rectangle, rad int, c1, c2 color.RGBA) {
	x0, w := r.Min.X, float64(r.Dx())
	fillRoundRectFunc(dst, r, rad, func(x int) color.RGBA {
		return lerp(c1, c2, float64(x-x0)/w)
	})
}

func strokeRoundRect(dst *image.RGBA, r image.Rectangle, rad int, c color.RGBA, w float64) {
	inner := image.Rect(r.Min.X+int(w), r.Min.Y+int(w), r.Max.X-int(w), r.Max.Y-int(w))
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			out := inRoundRect(float64(x)+0.5, float64(y)+0.5, r, float64(rad))
			in := inRoundRect(float64(x)+0.5, float64(y)+0.5, inner, float64(rad)-w)
			cov := out - in // ring coverage
			if cov > 0 {
				blendPixel(dst, x, y, c, cov)
			}
		}
	}
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
