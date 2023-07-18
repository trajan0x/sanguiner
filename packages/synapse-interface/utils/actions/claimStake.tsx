import { Address } from 'wagmi'
import toast from 'react-hot-toast'

import { MINICHEF_ADDRESSES } from '@/constants/minichef'
import ExplorerToastLink from '@/components/ExplorerToastLink'
import { txErrorHandler } from '@utils/txErrorHandler'
import { harvestLpPool } from '@/actions/harvestLpPool'

export const claimStake = async (
  chainId: number,
  address: Address,
  poolId: number
) => {
  let pendingPopup: any
  let successPopup: any

  pendingPopup = toast(`Starting your claim...`, {
    id: 'claim-in-progress-popup',
    duration: Infinity,
  })

  try {
    if (!address) throw new Error('Wallet must be connected')
    const tx = await harvestLpPool({
      address,
      chainId,
      poolId,
      lpAddress: MINICHEF_ADDRESSES[chainId],
    })

    toast.dismiss(pendingPopup)

    const successToastContent = (
      <div>
        <div>Claim Completed:</div>
        <ExplorerToastLink
          transactionHash={tx?.transactionHash}
          chainId={chainId}
        />
      </div>
    )

    successPopup = toast.success(successToastContent, {
      id: 'claim-success-popup',
      duration: 10000,
    })

    return tx
  } catch (err) {
    toast.dismiss(pendingPopup)
    txErrorHandler(err)
  }
}
