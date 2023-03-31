import { Provider } from '@ethersproject/abstract-provider'
import { providers as etherProvider } from 'ethers'
import { BigNumber } from '@ethersproject/bignumber'

import { SynapseSDK } from './sdk'

describe('SynapseSDK', () => {
  const arbitrumProvider: Provider = new etherProvider.JsonRpcProvider(
    'https://arb1.arbitrum.io/rpc'
  )
  const avalancheProvider: Provider = new etherProvider.JsonRpcProvider(
    'https://api.avax.network/ext/bc/C/rpc'
  )

  describe('#constructor', () => {
    it('fails with unequal amount of chains to providers', () => {
      const chainIds = [42161, 43114]
      const providers = [arbitrumProvider]
      expect(() => new SynapseSDK(chainIds, providers)).toThrowError(
        'Amount of chains and providers does not equal'
      )
    })
  })

  describe('bridgeQuote', () => {
    it('test', async () => {
      const chainIds = [42161, 43114]
      const providers = [arbitrumProvider, avalancheProvider]
      const Synapse = new SynapseSDK(chainIds, providers)
      const { bridgeFee, destQuery } = await Synapse.bridgeQuote(
        42161,
        43114,
        '0x8D9bA570D6cb60C7e3e0F31343Efe75AB8E65FB1',
        '0x321E7092a180BB43555132ec53AaA65a5bF84251',
        BigNumber.from('10000000000000000000')
      )
      expect(bridgeFee).toBeGreaterThan(0)
      expect(destQuery?.length).toBeGreaterThan(0)
      console.log(destQuery)
    })
  })

  describe('bridge', () => {
    it('test', async () => {
      const chainIds = [42161, 43114]
      const providers = [arbitrumProvider, avalancheProvider]
      const Synapse = new SynapseSDK(chainIds, providers)
      const { originQuery, destQuery } = await Synapse.bridgeQuote(
        42161,
        43114,
        '0x8D9bA570D6cb60C7e3e0F31343Efe75AB8E65FB1',
        '0x321E7092a180BB43555132ec53AaA65a5bF84251',
        BigNumber.from('10000000000000000000')
      )
      const { data, to } = await Synapse.bridge(
        '0x0AF91FA049A7e1894F480bFE5bBa20142C6c29a9',
        42161,
        43114,
        '0xff970a61a04b1ca14834a43f5de4533ebddb5cc8',
        BigNumber.from('20000000'),
        originQuery!,
        destQuery!
      )
      console.log(data, 'ttttt', to)
      expect(data?.length).toBeGreaterThan(0)
      expect(to?.length).toBeGreaterThan(0)
    })
  })
})
