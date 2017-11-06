Atlas Test-Net Example
======================

Atlas is a community test-net which can be used to test setting up a
cosmos validator node. Supplimentary to this tutorial you can also
follow `this video <https://www.youtube.com/watch?v=B-shjoqvnnY>`__.

To work on the Atlas you will need some tokens to get started. To do
this first generate a new key:

::

    MYNAME=<your name>
    gaiacli keys new $MYNAME
    gaiacli keys list
    MYADDR=<your newly generated address>

Then enter your key into `this
utility <http://www.cosmosvalidators.com/>`__ with your key address and
it will send you some ``fermion`` testnet tokens :). Fermions are the
native staking token on Atlas - see below for how to check your balance.

Now, to sync with the testnet, we need the genesis file and seeds. The
easiest way to get them is to Fetch and navigate to the tendermint
testnet repo:

::

    git clone https://github.com/tendermint/testnets $HOME/testnets
    GAIANET=$HOME/testnets/atlas/gaia
    cd $GAIANET

Now we can start a new node in the background - but note that it may
take a decent chunk of time to sync with the existing testnet... go brew
a pot of coffee!

::

    gaia start --home=$GAIANET  &> atlas.log &

Of course, you can follow the logs to see progress with
``tail -f atlas.log``. Once blocks slow down to about one block per
second, you're all caught up.

The ``gaia start`` command will automaticaly generate a validator
private key found in ``$GAIANET/priv_validator.json``. Let's get the
pubkey data for our validator node. The pubkey is located under
``"pub_key"{"data":`` within the json file

::

    cat $GAIANET/priv_validator.json 
    PUBKEY=<your newly generated address>  

If you have a json parser like ``jq``, you can get just the pubkey like
so:

::

    PUBKEY=$(cat $GAIANET/priv_validator.json | jq -r .pub_key.data)

Next let's initialize the gaia client to start interacting with the
testnet:

::

    gaiacli init --chain-id=atlas --node=tcp://localhost:46657

And check our balance:

::

    gaiacli query account $MYADDR

We are now ready to bond some tokens:

::

    gaiacli tx bond --amount=5fermion --name=$MYNAME --pubkey=$PUBKEY

Bonding tokens means that your balance is tied up as *stake*. Don't
worry, you'll be able to get it back later. As soon as some tokens have
been bonded the validator node which we started earlier will have power
in the network and begin to participate in consensus!

We can now check the validator set and see that we are a part of the
club!

::

    gaiacli query validators

Finally lets unbond to get back our tokens

::

    gaiacli tx unbond --amount=5fermion --name=$MYNAME

Remember to unbond before stopping your node!
