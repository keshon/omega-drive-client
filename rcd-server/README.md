# Server
Server (compiles to `rcd.exe`) allows front-end (which is `win-client`) to work with S3 operations (mount / unmount and etc.).

Servier is a wrapped version of [Rclone](https://github.com/rclone/rclone "Rclone"). The wrapping is needed to be sure that server will not work if front-end dies (because front-end is responsible for verifyng user access).

To compile server (Rclone) with mount capablities sources of [WinFsp](https://github.com/winfsp/winfsp "WinFsp") are required and MiniGW-64 must be installed.