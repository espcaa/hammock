package ui

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"sync"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

const maxRemoteImageBytes = 8 << 20

type ImageCache struct {
	mu      sync.Mutex
	items   map[string]image.Image
	loading map[string]bool
	client  *http.Client
}

func NewImageCache() *ImageCache {
	return &ImageCache{
		items:   make(map[string]image.Image),
		loading: make(map[string]bool),
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *ImageCache) Get(url string) (image.Image, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if img, ok := c.items[url]; ok {
		return img, true
	}
	if !c.loading[url] {
		c.loading[url] = true
		go c.load(url)
	}
	return nil, false
}

func (c *ImageCache) Pending(url string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.loading[url]
}

func (c *ImageCache) load(url string) {
	defer func() {
		c.mu.Lock()
		delete(c.loading, url)
		c.mu.Unlock()
	}()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 15_6_0)")

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return
	}

	img, _, err := image.Decode(io.LimitReader(resp.Body, maxRemoteImageBytes))
	if err != nil {
		return
	}

	c.mu.Lock()
	c.items[url] = img
	c.mu.Unlock()
}

type RemoteImage struct {
	cache       *ImageCache
	Size        unit.Dp
	Radius      unit.Dp
	URL         string
	Placeholder color.NRGBA
}

func (t *Theme) RemoteImage(cache *ImageCache, url string, size unit.Dp) RemoteImage {
	return RemoteImage{
		cache:       cache,
		Size:        size,
		Radius:      t.Radius,
		URL:         url,
		Placeholder: t.Overlay,
	}
}

func (r RemoteImage) Layout(gtx layout.Context) layout.Dimensions {
	sz := gtx.Dp(r.Size)
	rect := image.Rectangle{Max: image.Pt(sz, sz)}

	img, ok := r.cache.Get(r.URL)
	if ok {
		gtx.Constraints = layout.Exact(rect.Max)
		defer r.rrect(gtx, rect).Push(gtx.Ops).Pop()
		return widget.Image{Src: paint.NewImageOp(img), Fit: widget.Cover}.Layout(gtx)
	}

	gtx.Constraints = layout.Exact(rect.Max)
	defer r.rrect(gtx, rect).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, r.Placeholder)

	if r.cache.Pending(r.URL) {
		gtx.Execute(op.InvalidateCmd{})
	}
	return layout.Dimensions{Size: rect.Max}
}

func (r RemoteImage) rrect(gtx layout.Context, rect image.Rectangle) clip.Op {
	radius := gtx.Dp(r.Radius)
	return clip.RRect{Rect: rect, NW: radius, NE: radius, SW: radius, SE: radius}.Op(gtx.Ops)
}
