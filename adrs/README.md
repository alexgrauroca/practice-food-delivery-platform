# ADRs

## Adding a new ADR

To create a new ADR, use the following command:

```bash
make adr-new dir=<directory to create the adr> title="<adr title>"
```

For example, to create a new ADR for Go tech-stack:

```bash
make adr-new dir=tech-stack/go title="The new Go ADR title"
```

All the ADRs are created by using the [template](./templates/template.md).

## Supersede an existing ADR

To supersede an existing ADR, use the following command:

```bash
make adr-new dir=<directory to create the adr> title="<adr title>" super=<adr number to be superseded>
```

For example, to supersede the first ADr in Go tech-stack:

```bash
make adr-new dir=tech-stack/go title="The second ADR title" super=1
```
