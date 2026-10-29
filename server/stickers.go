package main

// stickerPack is the built-in sticker pack. The ids must match the files in
// client-web/public/stickers/<id>.svg (a test checks this). Only the id ever
// travels over the socket; the image is a static asset.
var stickerPack = []string{
	"smile", "laugh", "love", "cry", "angry", "surprised", "sleepy", "thumbs_up",
}

var stickerSet = func() map[string]bool {
	m := make(map[string]bool, len(stickerPack))
	for _, id := range stickerPack {
		m[id] = true
	}
	return m
}()

func isSticker(id string) bool { return stickerSet[id] }
