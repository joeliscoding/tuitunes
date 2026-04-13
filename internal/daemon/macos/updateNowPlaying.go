package macos

import "os/exec"

func UpdateNowPlaying(title, artist, album string) error {
	// TODO: Path to the helper should be in global config, not hardcoded
	// Should probably also use a more robust method of communicating with the helper, via unix socket, instead of just running it with arguments.
	cmd := exec.Command("./internal/daemon/macos/nowplayinghelper/tuitunes-nowplayinghelper", title, artist, album)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
