import {ellipsizeString} from '@utils/ellipsizeString'

export function StyleAddress({ sourceInfo, limiter = 4 }) {
  if (sourceInfo.address) {
    return (
      <span
        // className={`${getNetworkTextHoverColor(
        //   sourceInfo.chainId
        // )} hover:underline `}
        // href={getAddressesUrl({
        //   address: sourceInfo.address,
        //   chainIdTo: sourceInfo.chainId,
        // })}
        onClick={(e) => e.stopPropagation()}
      >
        {ellipsizeString({ string: sourceInfo.address, limiter, isZeroX: true })}
      </span>
    )
  } else {
    return '--'
  }
}
