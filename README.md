# digitalocean network tests

## Usage

After setting pulumi and the `DIGITALOCEAN_TOKEN` env var, you can run the following commands to deploy the infrastructure:

```bash
pulumi up
```

To destroy the infrastructure, you can run the following command:

```bash
pulumi down
```

### Generate Genesis

Use the netgen command to generate the genesis file and keys for the validators. This stores these files in the `./payload` directory.

```bash
netgen <number of validators> <chain-id>
```

