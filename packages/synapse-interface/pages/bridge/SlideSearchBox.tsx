import { useRef, useEffect } from 'react'

export default function SlideSearchBox({ searchStr, onSearch, placeholder }) {
  const inputRef = useRef(null)
  useEffect(() => inputRef.current?.focus(), [])

  return (
    <input
      ref={inputRef}
      className={`
        text-white 
        focus:outline-none
        hidden sm:inline-block
        flex-grow
        h-full min-w-[70%]
        py-2 pr-2
        rounded
        bg-transparent
       placeholder-white placeholder-opacity-40
      `}
      placeholder={placeholder}
      onChange={(e) => onSearch(e.target.value)}
      value={searchStr}
    />
  )
}
