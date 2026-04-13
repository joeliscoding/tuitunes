/*
 TODO:
 - play/pause/skip commands from control center
    - needs to send updates back to the daemon so it can update the state
 - album art
 - playback position
 - playback state (playing/paused)
 - handle updates via unix socket
 - this script should probably only take the file path of the track as an argument, and then read 
   the metadata itself, instead of having the daemon pass the metadata as arguments.
*/

import Foundation
import MediaPlayer

let commandCenter = MPRemoteCommandCenter.shared()
commandCenter.playCommand.isEnabled = true
commandCenter.playCommand.addTarget { _ in .success }
commandCenter.pauseCommand.isEnabled = true
commandCenter.pauseCommand.addTarget { _ in .success }

func updateInfo(title: String, artist: String) {
    let infoCenter = MPNowPlayingInfoCenter.default()
    var info = [String: Any]()
    info[MPMediaItemPropertyTitle] = title
    info[MPMediaItemPropertyArtist] = artist
    info[MPNowPlayingInfoPropertyPlaybackRate] = 1.0
    infoCenter.nowPlayingInfo = info
}

if CommandLine.arguments.count >= 3 {
    updateInfo(title: CommandLine.arguments[1], artist: CommandLine.arguments[2])
}

RunLoop.current.run()
