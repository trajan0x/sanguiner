import Link from 'next/link'
import { getPoolUrl } from '@urls'
import { switchNetwork } from '@wagmi/core'
import { useEffect, useState } from 'react'
import { getPoolData } from '@utils/actions/getPoolData'
import { getPoolApyData } from '@utils/actions/getPoolApyData'
import { Token } from '@types'
import Card from '@tw/Card'
import Grid from '@tw/Grid'
import { memo } from 'react'
import { CHAINS_BY_ID } from '@constants/chains'
import LoadingSpinner from '@tw/LoadingSpinner'

const PoolsListCard = memo(
  ({
    pool,
    chainId,
    address,
    connectedChainId,
    prices,
  }: {
    pool: Token
    chainId: number
    address: string
    connectedChainId: number
    prices: any
  }) => {
    const [poolData, setPoolData] = useState(undefined)
    const [poolApyData, setPoolApyData] = useState(undefined)

    console.log('PoolsListCard RERENDER')
    useEffect(() => {
      if (connectedChainId && chainId && address && pool) {
        // TODO - separate the apy and tvl so they load async.
        getPoolData(chainId, pool, address, false, prices)
          .then((res) => {
            console.log('res getPoolData', res)

            setPoolData(res)
          })
          .catch((err) => {
            console.log('Could not get Pool Data: ', err)
          })
        getPoolApyData(chainId, pool, prices)
          .then((res) => {
            console.log('res', res)
            setPoolApyData(res)
          })
          .catch((err) => {
            console.log('Could not get Pool APY Data: ', err)
          })
      }
    }, [])
    const chain = CHAINS_BY_ID[chainId]
    // const poolRouterIndex = POOL_INVERTED_ROUTER_INDEX[chainId][poolName]

    return (
      <div>
        <Link
          onClick={() => {
            if (address === undefined) {
              return alert('Please connect your wallet')
            }
            if (chainId != connectedChainId) {
              const res = switchNetwork({ chainId: chainId })
                .then((res) => {
                  return res
                })
                .catch(() => {
                  return undefined
                })
              if (res === undefined) {
                console.log("can't switch chain, chainId: ", chainId)
                return
              }
              // history.push(getPoolUrl({ poolRouterIndex }))
            }
          }}
          // href={getPoolUrl({ poolRouterIndex })}
          href={getPoolUrl(pool)} // TODO: fix this
        >
          <Card
            title={
              <PoolsCardTitle
                chainImg={chain?.chainImg?.src}
                poolName={pool.poolName}
                chainName={chain?.name}
              />
            }
            titleClassName="text-white font-light text-xl"
            className={`
            bg-bgBase transition-all rounded-xl items-center
            hover:bg-bgLight
            py-6 mt-4 pr-2
            border border-transparent
          `}
            divider={false}
          >
            <Grid gap={3} cols={{ xs: 3 }} className="mt-8">
              <div>
                <h3 className="text-sm text-opacity-50 text-secondaryTextColor">
                  Assets
                </h3>
                {pool.poolTokens && (
                  <CoinLabels coins={pool.poolTokens} /> // change coin to token
                )}
              </div>
              <div>
                <h3 className="text-sm text-opacity-50 text-secondaryTextColor">
                  TVL
                </h3>
                <div className={'mt-2 text-white '}>
                  {poolData?.totalLockedUSDStr ? (
                    '$' + poolData?.totalLockedUSDStr
                  ) : (
                    // <div className="animate-pulse rounded bg-slate-700 h-6 w-12" />
                    <LoadingSpinner />
                  )}
                </div>
              </div>
              <div>
                <h3 className="text-sm text-opacity-50 text-secondaryTextColor">
                  APY{' '}
                </h3>
                <div className="mt-2 text-green-400">
                  <>
                    {poolApyData?.fullCompoundedAPYStr
                      ? String(poolApyData.fullCompoundedAPYStr)
                      : '-'}
                    %
                  </>
                </div>
              </div>
            </Grid>
          </Card>
        </Link>
      </div>
    )
  }
)

const PoolsCardTitle = ({
  chainName,
  poolName,
  chainImg,
}: {
  chainName: string
  poolName: string
  chainImg: string
}) => {
  let displayPoolName = poolName?.replace(chainName, `<b>${chainName}</b>`)

  return (
    <div className="flex items-center">
      <img src={chainImg} className="w-6 h-6 mr-2 rounded-full" />
      {/* TODO: A better way to do this? */}
      <div dangerouslySetInnerHTML={{ __html: displayPoolName }} />
    </div>
  )
}

const CoinLabels = ({ coins }) => {
  return (
    <div className="flex mt-3">
      {coins
        ? coins?.map((coin, i) => (
            <img
              alt=""
              className="w-5 mr-1 rounded-full"
              src={coin.icon.src}
              key={i}
            />
          ))
        : null}
    </div>
  )
}
export default PoolsListCard
