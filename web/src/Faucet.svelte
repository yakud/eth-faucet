<script>
  import { onMount } from 'svelte';
  import { getAddress } from '@ethersproject/address';
  import { CloudflareProvider } from '@ethersproject/providers';
  import { setDefaults as setToast, toast } from 'bulma-toast';
  import detectEthereumProvider from '@metamask/detect-provider';

  let input = null;
  let faucetInfo = {
    account: '0x0000000000000000000000000000000000000000',
    network: 'testnet',
    payout: 1,
    tokenSymbol: 'ETH',
  };

  // try to load from local storage
  const localFaucetInfo = localStorage.getItem('faucetInfo');
  if (localFaucetInfo) {
    faucetInfo = JSON.parse(localFaucetInfo);
  }

  $: document.title = `${faucetInfo.tokenSymbol} ${capitalize(faucetInfo.network)} Faucet`;

  onMount(async () => {
    const res = await fetch('/api/info');
    faucetInfo = await res.json();

    // put the faucet info into the local storage
    localStorage.setItem('faucetInfo', JSON.stringify(faucetInfo));
  });

  setToast({
    position: 'bottom-center',
    dismissible: true,
    pauseOnHover: true,
    closeOnClick: false,
    animate: { in: 'fadeIn', out: 'fadeOut' },
  });

  async function handleRequest() {
    let address = input;
    if (address.endsWith('.eth')) {
      try {
        const provider = new CloudflareProvider();
        address = await provider.resolveName(address);
        if (!address) {
          toast({ message: 'invalid ENS name', type: 'is-warning' });
          return;
        }
      } catch (error) {
        toast({ message: error.reason, type: 'is-warning' });
        return;
      }
    }

    try {
      address = getAddress(address);
    } catch (error) {
      toast({ message: error.reason, type: 'is-warning' });
      return;
    }

    const res = await fetch('/api/claim', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        address,
      }),
    });

    let { msg } = await res.json();
    let type = res.ok ? 'is-success' : 'is-warning';
    toast({ message: msg, type });
  }

  function capitalize(str) {
    const lower = str.toLowerCase();
    return str.charAt(0).toUpperCase() + lower.slice(1);
  }

  let ethProvider;
  let accountAddress = '';

  onMount(async () => {
    ethProvider = await detectEthereumProvider();
    ethProvider.addListener('accountsChanged', (accounts) => {
      if (accounts.length === 0) {
        window.location.reload();
        accountAddress = '';
        return;
      }
      accountAddress = accounts[0];
    });

    ethProvider.addListener('chainChanged', (chainId) => {
      window.location.reload();
      accountAddress = '';
    });

    ethProvider.addListener('disconnect', (error) => {
      window.location.reload();
      accountAddress = '';
    });

    ethProvider.addListener('networkChanged', (networkId) => {
      window.location.reload();
      accountAddress = '';
    });
  });

  async function connectMetaMask() {
    if (!ethProvider) {
      alert('MetaMask is not installed');
      return;
    }

    try {
      const accounts = await ethProvider.request({ method: 'eth_requestAccounts' });
      accountAddress = accounts[0];
    } catch (error) {
      console.error('Error connecting to MetaMask:', error);
    }
  }
</script>

<main>
  <section class="hero is-info is-fullheight">
    <div class="hero-head">
      <nav class="navbar">
        <div class="container">
          <div class="navbar-brand">
            <a class="navbar-item" href="../..">
              <span class="icon">
                <i class="fa fa-bath" />
              </span>
              <span><b>{faucetInfo.network}</b></span>
            </a>
          </div>
          <div id="navbarMenu" class="navbar-menu">
            <div class="navbar-end">
              <span class="navbar-item">
                <a
                  class="button is-white is-outlined"
                  href="https://github.com/chainflag/eth-faucet"
                >
                  <span class="icon">
                    <i class="fa fa-github" />
                  </span>
                  <span>View Source</span>
                </a>
              </span>
            </div>
          </div>
        </div>
      </nav>
    </div>

    <div class="hero-body">
      <div class="container has-text-centered">
        <div class="column is-6 is-offset-3">
          <h1 class="title">
            Receive {faucetInfo.payout} {faucetInfo.tokenSymbol} per request
          </h1>
          <h2 class="subtitle">
            Serving from {faucetInfo.account}
          </h2>
          <div class="box">
            <div class="field is-grouped">
              <p class="control is-expanded has-icons-right">
                <input
                  value={accountAddress}
                  class="input is-rounded"
                  type="text"
                  placeholder="Enter your address or ENS name"
                />
                <button class="icon-button" type="button" on:click={connectMetaMask}>
                  <span class="icon is-small is-clickable">
                    <img src="./metamask2.png" style="width: 20px; height: 20px;" />
                  </span>
                </button>
              </p>

              <p class="control">
                <button
                  on:click={handleRequest}
                  class="button is-primary is-rounded"
                >
                  Request
                </button>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</main>

<style>
  .hero.is-info {
    background: linear-gradient(rgb(0 0 0 / 58%), rgb(0 0 0 / 94%));
    -webkit-background-size: cover;
    -moz-background-size: cover;
    -o-background-size: cover;
    background-size: cover;
  }
  .hero .subtitle {
    padding: 3rem 0;
    line-height: 1.5;
  }
  .box {
    border-radius: 19px;
  }
  .icon-button {
    background: none;
    border: none;
    padding: 0;
    position: absolute;
    right: 45px;
    top: 15%;
    transform: translateY(-50%);
    cursor: pointer;
  }
  .is-clickable {
    cursor: pointer;
  }

  .input.is-rounded + .is-clickable {
    position: absolute;
    right: 10px;
    top: 50%;
    transform: translateY(-50%);
  }

</style>
