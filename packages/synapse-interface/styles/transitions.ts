export const COIN_SLIDE_OVER_PROPS = {
  appear: true,
  unmount: true,
  enter: 'transition duration-150 ease-out',
  enterFrom: 'transform-gpu scale-y-0 ',
  enterTo: 'transform-gpu scale-y-100 opacity-100',
  leave: 'transition duration-150 ease-out ',
  leaveFrom: 'transform-gpu scale-y-100 opacity-100',
  leaveTo: 'transform-gpu scale-y-0 opacity-50',
}

export const SECTION_TRANSITION_PROPS = {
  enter: 'transition duration-100 ease-out',
  enterFrom: 'transform-gpu scale-y-0 ',
  enterTo: 'transform-gpu scale-y-100 opacity-100',
  leave: 'transition duration-75 ease-out ',
  leaveFrom: 'transform-gpu scale-y-100 opacity-100',
  leaveTo: 'transform-gpu scale-y-0 ',
  className: 'origin-top -mx-0 md:-mx-6',
}
