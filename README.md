# RectangleWin

A minimalistic Windows rewrite of [macOS
Rectangle.app](https://rectangleapp.com).

## Why?

It seems that no window snapping utility for Windows is capable of letting
user snip windows to {left, right, top, bottom} {½, ⅔, ⅓} using configurable
shortcut keys, and center windows in a screen like Rectangle.app does, so I
wrote this small utility for myself.

## Features

- Edge snapping ({left, right, top, bottom} {½, ⅔, ⅓})
  - <kbd>Win</kbd> + <kbd>Alt</kbd> + <kbd>&larr;</kbd><kbd>&rarr;</kbd><kbd>&uarr;</kbd><kbd>&darr;</kbd>
  - Press multiple times to alternate between half, two-thirds and one-thirds.

- Maximize Window: <kbd>Win</kbd>+<kbd>Shift</kbd>+<kbd>F</kbd>

## Roadmap

- Centering a window on the display.

- System tray (+running without a cmd window i.e. `"-H=windowsgui"`)

- Run on startup.

- Corner snapping (i.e. snapping to {top,bottom}{left,right}{½, ⅔, ⅓}).

- **Multiple monitor support**: I don't need this right now and I don't own
  a secondary display so these will need your help.
  - Support multiple displays (the code is very likely buggy right now when the
  primary display isn't the leftmost-topmost in the display arrangement)
  - Moving a window between displays

- Configurable shortcuts: I don't need these and it will likely require a pop-up
  UI, so I will probably not get to this.

## License

See [LICENSE](./LICENSE).