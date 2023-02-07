import { QuestionMarkCircleIcon } from '@heroicons/react/outline'

import { CHAIN_EXPLORER_URLS, CHAIN_INFO_MAP } from '@constants/networks'
import { getNetworkTextColor } from '@styles/networks'
import Image from 'next/image'

export function ChainInfo({
  chainId,
  imgClassName = 'w-4 h-4',
  linkClassName = 'float-right text-white hover:bg-cyan-100/50 px-1.5  rounded-md ease-in-out bg-[#191919]',
  textClassName = "pl-1 whitespace-nowrap text-sm text-white",
  txHash,
  noLink = false,
}) {
  const { chainName, chainImg } = CHAIN_INFO_MAP[chainId] ?? {}
  const explorer = CHAIN_EXPLORER_URLS[chainId] ?? ''
  var link = explorer
  if (txHash) {
    link += '/tx/' + txHash
  }
  if (chainName) {
    return (
      <div className="flex justify-between">
        <div className='flex flex-row  items-center'>
          <Image
            className={`inline mr-1 rounded-lg ${imgClassName}`}
            src={chainImg}
            alt={chainImg}
          />
          <span className={textClassName}>{chainName}</span></div>
        {noLink ? null : (
          <div className='flex items-center'>
            <a
            type="link"
              target="_blank"
              href={link}
              className={linkClassName}
            >
              ↗
            </a></div>)}
      </div>
    )
  } else {
    return (
      <div className="flex items-center">
        <QuestionMarkCircleIcon
          className={`inline mr-2 rounded-lg ${imgClassName}`}
          strokeWidth={1}
        />
        <span>--</span>
      </div>
    )
  }
}
