![VaporTrace Logo](../../../assets/images/VaporTrace_Logo.png)

/*
Copyright (c) 2026 José María Micoli
Licensed under the Business Source License 1.1
Change Date: 2033-02-17
Change License: Apache-2.0

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

# Sprint 14: Cloud Pivoting & Multi-Cloud Exploitation

**Status:** ⏳ PLANNED | **Version:** v3.2-Hydra (Future) | **Target:** May-June 2026

---

## 🎯 Sprint Overview

Sprint 14 enables VaporTrace to pivot through cloud infrastructure, exploiting cloud-native services (AWS, Azure, GCP, Kubernetes). This sprint adds multi-cloud support, credential harvesting from cloud environments, and lateral movement within cloud ecosystems.

**Slogan:** "Cloud Infiltration - Multi-Cloud Pivoting"

---

## 📋 Planned Deliverables

### 14.1: Cloud Provider Abstraction Layer

**Status:** ⏳ PLANNED  
**Complexity:** High  
**Estimated Effort:** 100 hours

**Objective:** Unified interface for AWS, Azure, GCP, Kubernetes

**Design:**

1. **Cloud Provider Interface**
   - Uniform API for all cloud platforms
   - Provider-specific implementations
   - Automatic provider detection
   - Credential management per provider

2. **Supported Platforms**
   - AWS (EC2, S3, Lambda, RDS, IAM)
   - Microsoft Azure (VMs, Storage, App Service, SQL)
   - Google Cloud Platform (Compute Engine, Storage, App Engine)
   - Kubernetes (RBAC, Secrets, Service Accounts)
   - DigitalOcean (Droplets, Spaces, Kubernetes)

3. **Unified Operations**
   - ListResources() - Enumerate all resources
   - GetCredentials() - Extract credentials
   - ExecuteOnResource() - Run commands
   - EnumeratePermissions() - Check access levels
   - Pivot() - Move laterally

**Implementation Plan:**
```go
// Cloud Provider Interface
type CloudProvider interface {
    // Authenticate with cloud account
    Authenticate(creds CloudCredentials) error
    
    // List resources by type
    ListResources(ctx context.Context, resourceType string) ([]CloudResource, error)
    
    // Get resource details
    GetResource(ctx context.Context, resourceID string) (*CloudResource, error)
    
    // Execute code on resource
    ExecuteOnResource(ctx context.Context, resourceID, code string) (string, error)
    
    // Extract credentials from environment
    ExtractCredentials(ctx context.Context) (CloudCredentials, error)
    
    // Check current permissions
    GetPermissions(ctx context.Context) ([]Permission, error)
    
    // Enumerate other resources
    EnumerateNetwork(ctx context.Context) ([]NetworkResource, error)
}

// AWS Implementation
type AWSProvider struct {
    client *aws.Client
    region string
}

func (aws *AWSProvider) ListResources(ctx context.Context, resourceType string) ([]CloudResource, error) {
    // AWS-specific enumeration
    // Handle all resource types (EC2, RDS, S3, Lambda, etc.)
}

// Azure Implementation
type AzureProvider struct {
    client *azure.Client
}

// GCP Implementation
type GCPProvider struct {
    client *gcp.Client
}

// Kubernetes Implementation
type K8sProvider struct {
    client kubernetes.Interface
}

// Factory function
func DetectAndConnect(ctx context.Context) (CloudProvider, error) {
    // Auto-detect cloud platform from environment
    // AWS: Check for AWS_* env vars, EC2 metadata
    // Azure: Check for AZURE_* env vars
    // GCP: Check for GOOGLE_* env vars, GCP metadata
    // K8s: Check for kubeconfig, in-cluster secrets
}
```

**Status:** Planning phase

---

### 14.2: Cloud Credential Harvesting

**Status:** ⏳ PLANNED  
**Complexity:** Very High  
**Estimated Effort:** 120 hours

**Objective:** Extract and enumerate cloud credentials from compromised systems

**Credential Types:**

1. **AWS Credentials**
   - Access Keys (ID + Secret)
   - Temporary STS tokens
   - EC2 IAM instance profiles
   - S3 bucket credentials
   - Lambda environment variables
   - CloudFormation templates
   - Secrets Manager entries
   - Parameter Store entries

2. **Azure Credentials**
   - Application credentials
   - Service Principal secrets
   - Managed identities
   - Azure CLI cached tokens
   - PowerShell credential objects
   - Azure DevOps PATs
   - Key Vault secrets

3. **GCP Credentials**
   - Service account keys
   - OAuth tokens
   - GCP CLI credentials
   - Workload identity federation
   - Application Default Credentials (ADC)
   - Cloud Run service accounts

4. **Kubernetes Credentials**
   - Service account tokens
   - kubeconfig files
   - In-cluster mounted secrets
   - RBAC role/rolebinding enumeration

**Implementation Plan:**
```go
type CloudCredentialHarvester struct {
    providers map[string]CloudProvider
}

func (ch *CloudCredentialHarvester) HarvestAllCredentials() ([]CloudCredentials, error) {
    var creds []CloudCredentials
    
    // AWS credential harvesting
    awsCreds := ch.harvestAWSCredentials()
    creds = append(creds, awsCreds...)
    
    // Azure credential harvesting
    azureCreds := ch.harvestAzureCredentials()
    creds = append(creds, azureCreds...)
    
    // GCP credential harvesting
    gcpCreds := ch.harvestGCPCredentials()
    creds = append(creds, gcpCreds...)
    
    // Kubernetes credential harvesting
    k8sCreds := ch.harvestK8sCredentials()
    creds = append(creds, k8sCreds...)
    
    return creds, nil
}

func (ch *CloudCredentialHarvester) harvestAWSCredentials() []CloudCredentials {
    var creds []CloudCredentials
    
    // Check environment variables
    if akID := os.Getenv("AWS_ACCESS_KEY_ID"); akID != "" {
        creds = append(creds, CloudCredentials{
            Provider: "AWS",
            Type: "AccessKey",
            AccessKeyID: akID,
            SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
        })
    }
    
    // Check ~/.aws/credentials
    awsCredsFile := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")
    if creds, err := parseAWSCredsFile(awsCredsFile); err == nil {
        creds = append(creds, creds...)
    }
    
    // Check EC2 instance metadata (if running on EC2)
    if instanceCreds := fetchEC2InstanceCreds(); instanceCreds != nil {
        creds = append(creds, instanceCreds)
    }
    
    // Check Lambda environment (if running in Lambda)
    if lambdaCreds := fetchLambdaEnvironment(); lambdaCreds != nil {
        creds = append(creds, lambdaCreds)
    }
    
    return creds
}

func (ch *CloudCredentialHarvester) ValidateCredentials(creds CloudCredentials) bool {
    // Test credential validity
    // AWS: Attempt sts:GetCallerIdentity
    // Azure: Attempt list subscriptions
    // GCP: Attempt list projects
    // K8s: Attempt api-resources
}
```

**Status:** Planning phase

---

### 14.3: Cloud-Native Exploitation Modules

**Status:** ⏳ PLANNED  
**Complexity:** Very High  
**Estimated Effort:** 150 hours

**Objective:** Exploit cloud-specific vulnerabilities

**Cloud Exploitation Modules:**

1. **AWS Modules**
   - S3 Bucket Enumeration & Breaches
   - EC2 Security Group Misconfiguration
   - IAM Policy Exploitation (privilege escalation)
   - Lambda Function Enumeration & Invocation
   - RDS Database Access
   - CloudFront/CDN misconfiguration
   - Secrets Manager / Parameter Store extraction
   - Metadata Service (IMDSv1/v2) exploitation

2. **Azure Modules**
   - Subscription enumeration
   - Storage Account Enumeration
   - Virtual Machine RDP/SSH access
   - SQL Server enumeration
   - Key Vault access
   - App Service environment variable extraction
   - Role assignment enumeration
   - Managed identity exploitation

3. **GCP Modules**
   - Project enumeration
   - Compute Engine instance exploitation
   - Cloud Storage bucket enumeration
   - CloudSQL database access
   - Cloud Run service exploitation
   - Pub/Sub topic enumeration
   - Service account key extraction

4. **Kubernetes Modules**
   - Deployment enumeration
   - Pod execution (kubectl exec equivalent)
   - Secret extraction
   - Service account privilege escalation
   - RBAC policy enumeration
   - Persistent volume enumeration

**Implementation Plan:**
```go
type CloudExploitModule interface {
    Name() string
    Provider() string
    Exploit(ctx context.Context, config ExploitConfig) (*ExploitResult, error)
    RequiredPermissions() []string
    Cleanup() error
}

// AWS S3 Bucket Enumeration
type AWSS3BucketModule struct{}

func (m *AWSS3BucketModule) Exploit(ctx context.Context, config ExploitConfig) (*ExploitResult, error) {
    // 1. Enumerate all S3 buckets
    // 2. Check bucket policies
    // 3. List bucket contents
    // 4. Check for public access
    // 5. Download sensitive files
    // 6. Return findings
}

// Azure Storage Account Exploitation
type AzureStorageModule struct{}

func (m *AzureStorageModule) Exploit(ctx context.Context, config ExploitConfig) (*ExploitResult, error) {
    // 1. Enumerate storage accounts
    // 2. Check access keys
    // 3. List containers
    // 4. Download blobs
}

// Kubernetes Pod Execution
type K8sPodExecModule struct{}

func (m *K8sPodExecModule) Exploit(ctx context.Context, config ExploitConfig) (*ExploitResult, error) {
    // 1. Find target pod
    // 2. Establish exec session
    // 3. Execute commands
    // 4. Return output
}
```

**Status:** Planning phase

---

### 14.4: Cross-Cloud Lateral Movement

**Status:** ⏳ PLANNED  
**Complexity:** Very High  
**Estimated Effort:** 100 hours

**Objective:** Move between cloud providers and on-premises systems

**Design:**

1. **Lateral Movement Chains**
   - AWS → Azure: Shared credentials
   - Azure → GCP: Cross-tenant migration
   - GCP → Kubernetes: Service account exploitation
   - Kubernetes → AWS: Cloud provider integration
   - Any → On-Premises: VPN/Direct Connect

2. **Credential Reuse**
   - AWS keys used for Azure (if configured)
   - Service account tokens across environments
   - SSH keys stored in cloud systems
   - Shared password managers (AWS Secrets, Azure Key Vault)

3. **Network-Based Pivoting**
   - VPN/Direct Connect traversal
   - Kubernetes service mesh exploitation
   - Cross-region VPC peering
   - Transit Gateway abuse

**Implementation Plan:**
```go
type CloudPivot struct {
    currentProvider CloudProvider
    discoveredCreds []CloudCredentials
}

func (cp *CloudPivot) EnumerateCrossCloudAccess() map[string]CloudProvider {
    // Find credentials for other cloud providers
    // Attempt to authenticate
    // Return accessible providers
}

func (cp *CloudPivot) PivotToProvider(provider string) error {
    // Switch to new provider
    // Establish connection
    // Update current context
}

func (cp *CloudPivot) BuildLateralMovementChain() (*MovementChain, error) {
    // 1. Identify current location
    // 2. Enumerate accessible systems
    // 3. Find paths to sensitive targets
    // 4. Generate movement chain
    // 5. Execute sequentially
}
```

**Status:** Planning phase

---

## 🔄 Current Status by Sub-Phase

| Sub-Phase | Deliverable | Status | Est. Completion |
|-----------|-------------|--------|-----------------|
| **14.1** | Cloud Abstraction | ⏳ PLANNED | June 2026 |
| **14.2** | Credential Harvesting | ⏳ PLANNED | June 2026 |
| **14.3** | Cloud Exploitation | ⏳ PLANNED | July 2026 |
| **14.4** | Cross-Cloud Movement | ⏳ PLANNED | July 2026 |

**Overall Progress:** 0% (Planning phase)

---

## 🏗️ Architecture

### Cloud Exploitation Stack

```
┌──────────────────────────────────────────────────┐
│ Operator / ProcessChain()                        │
├──────────────────────────────────────────────────┤
│ Cloud Pivot Controller                           │
├──────────────────────────────────────────────────┤
│ Cloud Provider Abstraction Layer                 │
├─────────────────┬──────────────┬─────────────────┤
│  AWS Provider   │ Azure Prov.  │  GCP Provider   │
├─────────────────┼──────────────┼─────────────────┤
│ ├─ EC2          │ ├─ VMs       │ ├─ Compute Eng. │
│ ├─ S3           │ ├─ Storage   │ ├─ Storage      │
│ ├─ Lambda       │ ├─ App Srv.  │ ├─ Run          │
│ ├─ RDS          │ ├─ SQL       │ ├─ SQL          │
│ └─ Secrets      │ └─ Key Vault │ └─ Secrets      │
└─────────────────┴──────────────┴─────────────────┘
        ↓              ↓                  ↓
    AWS APIs       Azure APIs         GCP APIs
    (HTTPS)        (HTTPS)           (HTTPS)
```

### Credential Harvesting Flow

```
1. Detect cloud environment
2. Scan for credentials:
   - Environment variables
   - Config files
   - Metadata services
   - Application memory
   - Kubernetes mounts
3. Validate credentials
4. Store securely in vault
5. Attempt lateral movement
```

---

## 📊 Expected Capabilities

| Feature | V3.1-Hydra | V3.2-Hydra (14+) |
|---------|-----------|-----------------|
| **AWS Exploitation** | Limited | ✅ FULL |
| **Azure Exploitation** | ❌ | ✅ NEW |
| **GCP Exploitation** | ❌ | ✅ NEW |
| **Kubernetes Exploitation** | ❌ | ✅ NEW |
| **Cloud Credential Harvesting** | ❌ | ✅ NEW |
| **Cross-Cloud Pivoting** | ❌ | ✅ NEW |
| **Multi-Cloud Orchestration** | ❌ | ✅ NEW |
| **Cloud API Abuse** | Limited | ✅ FULL |

---

## 🎯 Success Criteria

### Sprint 14 Completion Requirements

- [ ] Cloud Provider abstraction layer complete
- [ ] All 4 major cloud platforms supported
- [ ] Credential harvesting for all platforms
- [ ] 20+ cloud-specific exploitation modules
- [ ] Cross-cloud lateral movement functional
- [ ] Integration with ProcessChain() complete
- [ ] Full documentation and examples
- [ ] Security audit completed

---

## 📚 Documentation References

- **Dev-Roadmap.md** - Full roadmap context
- **Sprint-13/README.md** - C2 Architecture foundation
- **Sprint-11/README.md** - Autonomy engine
- **Sprint-16/README.md** - Verification patterns

---

## 🚀 Integration Points

### With Sprint 13 (C2)
- C2 agents deployed to cloud instances
- Cloud resources managed via C2
- Distributed exploitation across cloud

### With Sprint 15 (Mastery)
- Advanced exploitation chains
- Multi-cloud orchestration templates
- Lateral movement automation

### With ProcessChain()
- Cloud-aware chain execution
- Dynamic credential injection
- Cross-provider chain steps

---

## ⚠️ Security & Compliance Challenges

### Challenge 1: Cloud IAM Complexity
- **Problem:** Each cloud has unique IAM model
- **Mitigation:** Abstraction layer, standardized permission checking
- **Status:** Core design of 14.1

### Challenge 2: Credential Exposure
- **Problem:** Storing cloud credentials locally
- **Mitigation:** Encrypted vault, in-memory only, auto-rotation
- **Status:** Integration with Sprint 16 verification

### Challenge 3: Service Limits
- **Problem:** Cloud API rate limiting
- **Mitigation:** Throttling, batch operations, caching
- **Status:** Design phase

### Challenge 4: Audit Logging
- **Problem:** Actions appear in cloud provider logs
- **Mitigation:** Use legitimate credentials, batch high-volume operations
- **Status:** Operational concern, depends on engagement rules

---

## 📊 Estimated Metrics

| Metric | Target |
|--------|--------|
| **Supported Cloud Platforms** | 4+ |
| **Exploitation Modules** | 20+ |
| **Credential Types Harvested** | 15+ |
| **Lateral Movement Chains** | 10+ |
| **API Latency** | <2 seconds |
| **Success Rate (Exploitation)** | >85% |

---

## 📄 Status Summary

**Sprint 14 Status:** ⏳ PLANNED (Not yet started)

**Target Timeline:** May-June 2026  
**Estimated Effort:** 470 hours (12 weeks)  
**Dependencies:** Sprint 13 completion (C2 Architecture)

**Impact:** Enables VaporTrace to pivot through multi-cloud environments, dramatically expanding scope in enterprise penetration tests.

---

**Status:** ⏳ PLANNED | **Last Updated:** February 8, 2026 | **Next Review:** March 15, 2026
