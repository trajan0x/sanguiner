import Head from 'next/head'
import BridgePage from './bridge'
import { LandingPageWrapper } from '@/components/layouts/LandingPageWrapper'
export default function Home() {
  return (
    <>
      <Head>
        <title>Synapse</title>
        <meta name="description" content="Generated by create next app" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <BridgePage />
    </>
  )
}
