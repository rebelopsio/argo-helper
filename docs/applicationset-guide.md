# ApplicationSet Guide

ApplicationSets enable declarative management of multiple ArgoCD Applications through a single resource. This guide explains the different types of ApplicationSet generators and when to use them.

## Generator Types

### Git Generator

Use the Git generator when you want to create Applications based on the structure of a Git repository.

```yaml
generators:
  - git:
      repoURL: https://github.com/yourorg/your-repo.git
      revision: HEAD
      directories:
        - path: apps/*
```

**Best for:**
- Mono-repo structures where applications are organized in directories
- When you want to add new applications by just creating a new directory

### Cluster Generator

Use the Cluster generator to deploy the same application across multiple clusters.

```yaml
generators:
  - clusters:
      selector:
        matchLabels:
          environment: production
```

**Best for:**
- Multi-cluster deployments
- Edge computing scenarios
- Consistent application rollout across environments

### List Generator

Use the List generator when you have a specific, predefined set of applications to deploy.

```yaml
generators:
  - list:
      elements:
        - name: app1
          namespace: app1-ns
          environment: dev
        - name: app2
          namespace: app2-ns
          environment: staging
```

**Best for:**
- Small, fixed sets of applications
- When applications need unique parameters
- Testing specific combinations

### Matrix Generator

Use the Matrix generator to combine two different generators, creating a Cartesian product of all combinations.

```yaml
generators:
  - matrix:
      generators:
        - git:
            repoURL: https://github.com/yourorg/your-repo.git
            revision: HEAD
            directories:
              - path: apps/*
        - clusters:
            selector:
              matchLabels:
                argocd.argoproj.io/secret-type: cluster
```

**Best for:**
- Deploying all applications to all environments
- Complex deployments requiring multiple dimensions
- When you need a full mesh of applications across infrastructure

### Merge Generator

Use the Merge generator to combine multiple generators using a merge key.

```yaml
generators:
  - merge:
      generators:
        - clusters: {}
        - list:
            elements:
              - app: frontend
                path: apps/frontend
      mergeKeys:
        - app
```

**Best for:**
- When you need to combine parameters from different sources
- More controlled deployments than Matrix
- Overlaying environment-specific configs on applications

## Best Practices

1. **Use Templating Carefully**: Use quotes around template values: `name: '{{name}}'`

2. **Test ApplicationSets First**: Use `--dry-run` to preview the Applications that will be created

3. **Progressive Rollouts**: Use the `ProgressiveSync` strategy for safer multi-cluster deployments

4. **Resource Organization**: Group related Applications in the same ApplicationSet

5. **Reuse Templates**: Define common template patterns as helpers

6. **Leverage Labels**: Use metadata labels for filtering and organization

7. **GitOps for ApplicationSets**: Store your ApplicationSets in Git for versioning and history

## Common Use Cases

- **Multi-Environment Deployments**: Using Cluster generator to deploy across dev/staging/prod
- **Microservices Deployment**: Using Git generator to deploy all services in a mono-repo
- **Multi-Tenant Applications**: Using Matrix to deploy tenant-specific configurations
- **Regional Deployments**: Deploying the same applications across geographic regions

## Example Implementations

See the `docs/examples/applicationset-templates/` directory for reference implementations of each generator type.