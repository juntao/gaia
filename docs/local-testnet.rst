Local-Test Example
==================

Here is a quick example to get you off your feet:

First, generate a new key with a name, and save the address:

::

    MYNAME=<your name>
    gaiacli keys new $MYNAME
    gaiacli keys list
    MYADDR=<your newly generated address>

Now initialize a gaia chain:

::

    gaia init $MYADDR --home=$HOME/.atlas1 --chain-id=test 

This will create all the files necessary to run a single node chain in
``$HOME/.atlas1``: a ``priv_validator.json`` file with the validators
private key, and a ``genesis.json`` file with the list of validators and
accounts. In this case, we have one random validator, and ``$MYADDR`` is
an independent account that has a bunch of coins.

We can add a second node on our local machine by initiating a node in a
new directory, and copying in the genesis:

::

    gaia init $MYADDR --home=$HOME/.atlas2 --chain-id=test
    cp $HOME/.atlas1/genesis.json $HOME/.atlas2/genesis.json

We need to also modify ``$HOME/.atlas2/config.toml`` to set new seeds
and ports. It should look like:

::

    proxy_app = "tcp://127.0.0.1:46668"
    moniker = "anonymous"
    fast_sync = true
    db_backend = "leveldb"
    log_level = "state:info,*:error"

    [rpc]
    laddr = "tcp://0.0.0.0:46667"

    [p2p]
    laddr = "tcp://0.0.0.0:46666"
    seeds = "0.0.0.0:46656"

Great, now that we've initialized the chains, we can start both nodes in
the background:

::

    gaia start --home=$HOME/.atlas1  &> atlas1.log &
    NODE1_PID=$!
    gaia start --home=$HOME/.atlas2  &> atlas2.log &
    NODE2_PID=$!

Note we save the ``PID`` so we can later kill the processes.

Of course, you can peak at your logs with ``tail atlas1.log``, or follow
them for a bit with ``tail -f atlas1.log``.

Now we can initialize a client for the first node, and look up our
account:

::

    gaiacli init --chain-id=test --node=tcp://localhost:46657
    gaiacli query account $MYADDR

Nice. We can also lookup the validator set:

::

    gaiacli query validators

Notice it's empty! This is because the initial validators are special -
the app doesn't know about them, so they can't be removed. To see what
tendermint itself thinks the validator set is, use:

::

    curl localhost:46657/validators

Ok, let's add the second node as a validator. First, we need the pubkey
data:

::

    cat $HOME/.atlas2/priv_validator.json 

If you have a json parser like ``jq``, you can get just the pubkey:

::

    cat $HOME/.atlas2/priv_validator.json | jq .pub_key.data

Now we can bond some coins to that pubkey:

::

    gaiacli tx bond --amount=10fermion --name=$MYNAME --pubkey=<validator pubkey>

We should see our account balance decrement, and the pubkey get added to
the app's list of bonds:

::

    gaiacli query account $MYADDR
    gaiacli query validators

To confirm for certain the new validator is active, check tendermint:

::

    curl localhost:46657/validators

If you now kill your second node, blocks will stop streaming in, because
there aren't enough validators online. Turn her back on and they will
start streaming again.

Finally, to relinquish all your power, unbond some coins. You should see
your VotingPower reduce and your account balance increase.

::

    gaiacli tx unbond --amount=10fermion --name=$MYNAME
    gaiacli query validators
    gaiacli query account $MYADDR

Once you unbond enough, you will no longer be needed to make new blocks.
