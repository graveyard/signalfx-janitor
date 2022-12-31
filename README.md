# signalfx-janitor

Cleans up stale alerts in SignalFX,
Alternatively mutes noisy alerts for a specified amount of time.

Owned by `eng-infra`.

## Running locally

Build and run the app:

```
export SFX_TOKEN=<signalfx api token>
export SFX_ORG_ID=<signalfx organization ID>
make run
```

or via ark:

```
ark start --local
```

## Deploying

```
ark start signalfx-janitor -e production
```
test
