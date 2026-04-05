// Copyright (c) 2026 TerraKit. Licensed under BSL 1.1.

package pricing

// HoursPerMonth is the standard cloud-billing assumption used by AWS, GCP, and
// Azure (730 h ≈ 365 days / 12 months).
const HoursPerMonth = 730.0

// SKU holds the pricing data for a single cloud resource variant.
//
// JSON wire format (used by local_pricing_path files):
//
//	{
//	  "hourly_usd":  0.0104,
//	  "description": "t3.micro On-Demand (us-east-1)"
//	}
type SKU struct {
	HourlyUSD   float64 `json:"hourly_usd"`
	Description string  `json:"description,omitempty"`
}

// PricingDB is the top-level type for both the built-in catalogue and any
// user-supplied JSON file passed via local_pricing_path.
//
// Structure:  resource_type → sku_key → SKU
//
// Example JSON file:
//
//	{
//	  "aws_instance": {
//	    "t3.micro": { "hourly_usd": 0.0104, "description": "t3.micro on-demand" },
//	    "t3.small": { "hourly_usd": 0.0208, "description": "t3.small on-demand" }
//	  },
//	  "aws_db_instance": {
//	    "db.t3.micro": { "hourly_usd": 0.017, "description": "RDS t3.micro MySQL" }
//	  }
//	}
type PricingDB map[string]map[string]SKU

// builtinPricingDB returns the embedded Hello-World pricing catalogue that is
// active when no local_pricing_path is configured.  These figures are
// illustrative on-demand prices for us-east-1 and should be replaced with a
// real database in production use.
func builtinPricingDB() PricingDB {
	return PricingDB{

		// =====================================================================
		//  A W S
		// =====================================================================

		// ── EC2 instances (on-demand Linux, us-east-1) ────────────────────────
		"aws_instance": {
			// T3 — burstable general-purpose
			"t3.nano":    {HourlyUSD: 0.0052, Description: "t3.nano On-Demand Linux (us-east-1)"},
			"t3.micro":   {HourlyUSD: 0.0104, Description: "t3.micro On-Demand Linux (us-east-1)"},
			"t3.small":   {HourlyUSD: 0.0208, Description: "t3.small On-Demand Linux (us-east-1)"},
			"t3.medium":  {HourlyUSD: 0.0416, Description: "t3.medium On-Demand Linux (us-east-1)"},
			"t3.large":   {HourlyUSD: 0.0832, Description: "t3.large On-Demand Linux (us-east-1)"},
			"t3.xlarge":  {HourlyUSD: 0.1664, Description: "t3.xlarge On-Demand Linux (us-east-1)"},
			"t3.2xlarge": {HourlyUSD: 0.3328, Description: "t3.2xlarge On-Demand Linux (us-east-1)"},
			// T3a — AMD burstable
			"t3a.nano":    {HourlyUSD: 0.0047, Description: "t3a.nano On-Demand Linux (us-east-1)"},
			"t3a.micro":   {HourlyUSD: 0.0094, Description: "t3a.micro On-Demand Linux (us-east-1)"},
			"t3a.small":   {HourlyUSD: 0.0188, Description: "t3a.small On-Demand Linux (us-east-1)"},
			"t3a.medium":  {HourlyUSD: 0.0376, Description: "t3a.medium On-Demand Linux (us-east-1)"},
			"t3a.large":   {HourlyUSD: 0.0752, Description: "t3a.large On-Demand Linux (us-east-1)"},
			"t3a.xlarge":  {HourlyUSD: 0.1504, Description: "t3a.xlarge On-Demand Linux (us-east-1)"},
			"t3a.2xlarge": {HourlyUSD: 0.3008, Description: "t3a.2xlarge On-Demand Linux (us-east-1)"},
			// M5 — general-purpose
			"m5.large":    {HourlyUSD: 0.096, Description: "m5.large On-Demand Linux (us-east-1)"},
			"m5.xlarge":   {HourlyUSD: 0.192, Description: "m5.xlarge On-Demand Linux (us-east-1)"},
			"m5.2xlarge":  {HourlyUSD: 0.384, Description: "m5.2xlarge On-Demand Linux (us-east-1)"},
			"m5.4xlarge":  {HourlyUSD: 0.768, Description: "m5.4xlarge On-Demand Linux (us-east-1)"},
			"m5.8xlarge":  {HourlyUSD: 1.536, Description: "m5.8xlarge On-Demand Linux (us-east-1)"},
			"m5.12xlarge": {HourlyUSD: 2.304, Description: "m5.12xlarge On-Demand Linux (us-east-1)"},
			"m5.16xlarge": {HourlyUSD: 3.072, Description: "m5.16xlarge On-Demand Linux (us-east-1)"},
			"m5.24xlarge": {HourlyUSD: 4.608, Description: "m5.24xlarge On-Demand Linux (us-east-1)"},
			// M6i — latest gen general-purpose (Intel)
			"m6i.large":    {HourlyUSD: 0.096, Description: "m6i.large On-Demand Linux (us-east-1)"},
			"m6i.xlarge":   {HourlyUSD: 0.192, Description: "m6i.xlarge On-Demand Linux (us-east-1)"},
			"m6i.2xlarge":  {HourlyUSD: 0.384, Description: "m6i.2xlarge On-Demand Linux (us-east-1)"},
			"m6i.4xlarge":  {HourlyUSD: 0.768, Description: "m6i.4xlarge On-Demand Linux (us-east-1)"},
			"m6i.8xlarge":  {HourlyUSD: 1.536, Description: "m6i.8xlarge On-Demand Linux (us-east-1)"},
			"m6i.12xlarge": {HourlyUSD: 2.304, Description: "m6i.12xlarge On-Demand Linux (us-east-1)"},
			"m6i.16xlarge": {HourlyUSD: 3.072, Description: "m6i.16xlarge On-Demand Linux (us-east-1)"},
			"m6i.24xlarge": {HourlyUSD: 4.608, Description: "m6i.24xlarge On-Demand Linux (us-east-1)"},
			// M7i — 7th gen general-purpose (Intel)
			"m7i.large":   {HourlyUSD: 0.1008, Description: "m7i.large On-Demand Linux (us-east-1)"},
			"m7i.xlarge":  {HourlyUSD: 0.2016, Description: "m7i.xlarge On-Demand Linux (us-east-1)"},
			"m7i.2xlarge": {HourlyUSD: 0.4032, Description: "m7i.2xlarge On-Demand Linux (us-east-1)"},
			"m7i.4xlarge": {HourlyUSD: 0.8064, Description: "m7i.4xlarge On-Demand Linux (us-east-1)"},
			// M6g — Graviton2 general-purpose
			"m6g.medium":  {HourlyUSD: 0.0385, Description: "m6g.medium On-Demand Linux (us-east-1)"},
			"m6g.large":   {HourlyUSD: 0.077, Description: "m6g.large On-Demand Linux (us-east-1)"},
			"m6g.xlarge":  {HourlyUSD: 0.154, Description: "m6g.xlarge On-Demand Linux (us-east-1)"},
			"m6g.2xlarge": {HourlyUSD: 0.308, Description: "m6g.2xlarge On-Demand Linux (us-east-1)"},
			"m6g.4xlarge": {HourlyUSD: 0.616, Description: "m6g.4xlarge On-Demand Linux (us-east-1)"},
			"m6g.8xlarge": {HourlyUSD: 1.232, Description: "m6g.8xlarge On-Demand Linux (us-east-1)"},
			// M7g — Graviton3 general-purpose
			"m7g.medium":  {HourlyUSD: 0.0408, Description: "m7g.medium On-Demand Linux (us-east-1)"},
			"m7g.large":   {HourlyUSD: 0.0816, Description: "m7g.large On-Demand Linux (us-east-1)"},
			"m7g.xlarge":  {HourlyUSD: 0.1632, Description: "m7g.xlarge On-Demand Linux (us-east-1)"},
			"m7g.2xlarge": {HourlyUSD: 0.3264, Description: "m7g.2xlarge On-Demand Linux (us-east-1)"},
			"m7g.4xlarge": {HourlyUSD: 0.6528, Description: "m7g.4xlarge On-Demand Linux (us-east-1)"},
			// C5 — compute-optimised
			"c5.large":    {HourlyUSD: 0.085, Description: "c5.large On-Demand Linux (us-east-1)"},
			"c5.xlarge":   {HourlyUSD: 0.170, Description: "c5.xlarge On-Demand Linux (us-east-1)"},
			"c5.2xlarge":  {HourlyUSD: 0.340, Description: "c5.2xlarge On-Demand Linux (us-east-1)"},
			"c5.4xlarge":  {HourlyUSD: 0.680, Description: "c5.4xlarge On-Demand Linux (us-east-1)"},
			"c5.9xlarge":  {HourlyUSD: 1.530, Description: "c5.9xlarge On-Demand Linux (us-east-1)"},
			"c5.12xlarge": {HourlyUSD: 2.040, Description: "c5.12xlarge On-Demand Linux (us-east-1)"},
			"c5.18xlarge": {HourlyUSD: 3.060, Description: "c5.18xlarge On-Demand Linux (us-east-1)"},
			// C6i — latest gen compute-optimised (Intel)
			"c6i.large":   {HourlyUSD: 0.085, Description: "c6i.large On-Demand Linux (us-east-1)"},
			"c6i.xlarge":  {HourlyUSD: 0.170, Description: "c6i.xlarge On-Demand Linux (us-east-1)"},
			"c6i.2xlarge": {HourlyUSD: 0.340, Description: "c6i.2xlarge On-Demand Linux (us-east-1)"},
			"c6i.4xlarge": {HourlyUSD: 0.680, Description: "c6i.4xlarge On-Demand Linux (us-east-1)"},
			"c6i.8xlarge": {HourlyUSD: 1.360, Description: "c6i.8xlarge On-Demand Linux (us-east-1)"},
			// C6g — Graviton2 compute-optimised
			"c6g.medium":  {HourlyUSD: 0.034, Description: "c6g.medium On-Demand Linux (us-east-1)"},
			"c6g.large":   {HourlyUSD: 0.068, Description: "c6g.large On-Demand Linux (us-east-1)"},
			"c6g.xlarge":  {HourlyUSD: 0.136, Description: "c6g.xlarge On-Demand Linux (us-east-1)"},
			"c6g.2xlarge": {HourlyUSD: 0.272, Description: "c6g.2xlarge On-Demand Linux (us-east-1)"},
			"c6g.4xlarge": {HourlyUSD: 0.544, Description: "c6g.4xlarge On-Demand Linux (us-east-1)"},
			// C7g — Graviton3 compute-optimised
			"c7g.medium":  {HourlyUSD: 0.0361, Description: "c7g.medium On-Demand Linux (us-east-1)"},
			"c7g.large":   {HourlyUSD: 0.0725, Description: "c7g.large On-Demand Linux (us-east-1)"},
			"c7g.xlarge":  {HourlyUSD: 0.145, Description: "c7g.xlarge On-Demand Linux (us-east-1)"},
			"c7g.2xlarge": {HourlyUSD: 0.290, Description: "c7g.2xlarge On-Demand Linux (us-east-1)"},
			// R5 — memory-optimised
			"r5.large":    {HourlyUSD: 0.126, Description: "r5.large On-Demand Linux (us-east-1)"},
			"r5.xlarge":   {HourlyUSD: 0.252, Description: "r5.xlarge On-Demand Linux (us-east-1)"},
			"r5.2xlarge":  {HourlyUSD: 0.504, Description: "r5.2xlarge On-Demand Linux (us-east-1)"},
			"r5.4xlarge":  {HourlyUSD: 1.008, Description: "r5.4xlarge On-Demand Linux (us-east-1)"},
			"r5.8xlarge":  {HourlyUSD: 2.016, Description: "r5.8xlarge On-Demand Linux (us-east-1)"},
			"r5.12xlarge": {HourlyUSD: 3.024, Description: "r5.12xlarge On-Demand Linux (us-east-1)"},
			// R6i — latest gen memory-optimised (Intel)
			"r6i.large":   {HourlyUSD: 0.126, Description: "r6i.large On-Demand Linux (us-east-1)"},
			"r6i.xlarge":  {HourlyUSD: 0.252, Description: "r6i.xlarge On-Demand Linux (us-east-1)"},
			"r6i.2xlarge": {HourlyUSD: 0.504, Description: "r6i.2xlarge On-Demand Linux (us-east-1)"},
			"r6i.4xlarge": {HourlyUSD: 1.008, Description: "r6i.4xlarge On-Demand Linux (us-east-1)"},
			// R6g — Graviton2 memory-optimised
			"r6g.medium":  {HourlyUSD: 0.0504, Description: "r6g.medium On-Demand Linux (us-east-1)"},
			"r6g.large":   {HourlyUSD: 0.1008, Description: "r6g.large On-Demand Linux (us-east-1)"},
			"r6g.xlarge":  {HourlyUSD: 0.2016, Description: "r6g.xlarge On-Demand Linux (us-east-1)"},
			"r6g.2xlarge": {HourlyUSD: 0.4032, Description: "r6g.2xlarge On-Demand Linux (us-east-1)"},
			// I3 — storage-optimised
			"i3.large":   {HourlyUSD: 0.156, Description: "i3.large On-Demand Linux (us-east-1)"},
			"i3.xlarge":  {HourlyUSD: 0.312, Description: "i3.xlarge On-Demand Linux (us-east-1)"},
			"i3.2xlarge": {HourlyUSD: 0.624, Description: "i3.2xlarge On-Demand Linux (us-east-1)"},
			"i3.4xlarge": {HourlyUSD: 1.248, Description: "i3.4xlarge On-Demand Linux (us-east-1)"},
			// GPU — P3/G4
			"p3.2xlarge":  {HourlyUSD: 3.06, Description: "p3.2xlarge On-Demand Linux (us-east-1)"},
			"p3.8xlarge":  {HourlyUSD: 12.24, Description: "p3.8xlarge On-Demand Linux (us-east-1)"},
			"p3.16xlarge": {HourlyUSD: 24.48, Description: "p3.16xlarge On-Demand Linux (us-east-1)"},
			"g4dn.xlarge":  {HourlyUSD: 0.526, Description: "g4dn.xlarge On-Demand Linux (us-east-1)"},
			"g4dn.2xlarge": {HourlyUSD: 0.752, Description: "g4dn.2xlarge On-Demand Linux (us-east-1)"},
			"g4dn.4xlarge": {HourlyUSD: 1.204, Description: "g4dn.4xlarge On-Demand Linux (us-east-1)"},
			"g4dn.8xlarge": {HourlyUSD: 2.176, Description: "g4dn.8xlarge On-Demand Linux (us-east-1)"},
			"g5.xlarge":    {HourlyUSD: 1.006, Description: "g5.xlarge On-Demand Linux (us-east-1)"},
			"g5.2xlarge":   {HourlyUSD: 1.212, Description: "g5.2xlarge On-Demand Linux (us-east-1)"},
			"g5.4xlarge":   {HourlyUSD: 1.624, Description: "g5.4xlarge On-Demand Linux (us-east-1)"},
		},

		// ── EKS clusters (flat rate) ──────────────────────────────────────────
		"aws_eks_cluster": {
			"standard": {HourlyUSD: 0.10, Description: "EKS cluster endpoint ($0.10/hr flat)"},
		},

		// ── EKS node groups (priced by instance type) ─────────────────────────
		"aws_eks_node_group": {
			"t3.micro":    {HourlyUSD: 0.0104, Description: "EKS node t3.micro"},
			"t3.small":    {HourlyUSD: 0.0208, Description: "EKS node t3.small"},
			"t3.medium":   {HourlyUSD: 0.0416, Description: "EKS node t3.medium"},
			"t3.large":    {HourlyUSD: 0.0832, Description: "EKS node t3.large"},
			"t3.xlarge":   {HourlyUSD: 0.1664, Description: "EKS node t3.xlarge"},
			"t3.2xlarge":  {HourlyUSD: 0.3328, Description: "EKS node t3.2xlarge"},
			"m5.large":    {HourlyUSD: 0.096, Description: "EKS node m5.large"},
			"m5.xlarge":   {HourlyUSD: 0.192, Description: "EKS node m5.xlarge"},
			"m5.2xlarge":  {HourlyUSD: 0.384, Description: "EKS node m5.2xlarge"},
			"m5.4xlarge":  {HourlyUSD: 0.768, Description: "EKS node m5.4xlarge"},
			"m6i.large":   {HourlyUSD: 0.096, Description: "EKS node m6i.large"},
			"m6i.xlarge":  {HourlyUSD: 0.192, Description: "EKS node m6i.xlarge"},
			"m6i.2xlarge": {HourlyUSD: 0.384, Description: "EKS node m6i.2xlarge"},
			"m6g.large":   {HourlyUSD: 0.077, Description: "EKS node m6g.large"},
			"m6g.xlarge":  {HourlyUSD: 0.154, Description: "EKS node m6g.xlarge"},
			"m6g.2xlarge": {HourlyUSD: 0.308, Description: "EKS node m6g.2xlarge"},
			"c5.large":    {HourlyUSD: 0.085, Description: "EKS node c5.large"},
			"c5.xlarge":   {HourlyUSD: 0.170, Description: "EKS node c5.xlarge"},
			"c5.2xlarge":  {HourlyUSD: 0.340, Description: "EKS node c5.2xlarge"},
			"c6i.large":   {HourlyUSD: 0.085, Description: "EKS node c6i.large"},
			"c6i.xlarge":  {HourlyUSD: 0.170, Description: "EKS node c6i.xlarge"},
			"r5.large":    {HourlyUSD: 0.126, Description: "EKS node r5.large"},
			"r5.xlarge":   {HourlyUSD: 0.252, Description: "EKS node r5.xlarge"},
			"r5.2xlarge":  {HourlyUSD: 0.504, Description: "EKS node r5.2xlarge"},
		},

		// ── EKS Fargate ───────────────────────────────────────────────────────
		// Charged per vCPU-hr ($0.04048) + per GB-hr ($0.004445). Priced as
		// a fixed per-pod estimate based on typical 0.5 vCPU / 1 GB pod.
		"aws_eks_fargate_profile": {
			"standard": {HourlyUSD: 0.02469, Description: "EKS Fargate (0.5 vCPU + 1GB profile)"},
		},

		// ── RDS instances ─────────────────────────────────────────────────────
		"aws_db_instance": {
			"db.t3.micro":    {HourlyUSD: 0.017, Description: "RDS db.t3.micro MySQL"},
			"db.t3.small":    {HourlyUSD: 0.034, Description: "RDS db.t3.small MySQL"},
			"db.t3.medium":   {HourlyUSD: 0.068, Description: "RDS db.t3.medium MySQL"},
			"db.t3.large":    {HourlyUSD: 0.136, Description: "RDS db.t3.large MySQL"},
			"db.t4g.micro":   {HourlyUSD: 0.016, Description: "RDS db.t4g.micro MySQL"},
			"db.t4g.small":   {HourlyUSD: 0.032, Description: "RDS db.t4g.small MySQL"},
			"db.t4g.medium":  {HourlyUSD: 0.065, Description: "RDS db.t4g.medium MySQL"},
			"db.t4g.large":   {HourlyUSD: 0.129, Description: "RDS db.t4g.large MySQL"},
			"db.m5.large":    {HourlyUSD: 0.171, Description: "RDS db.m5.large MySQL"},
			"db.m5.xlarge":   {HourlyUSD: 0.342, Description: "RDS db.m5.xlarge MySQL"},
			"db.m5.2xlarge":  {HourlyUSD: 0.684, Description: "RDS db.m5.2xlarge MySQL"},
			"db.m5.4xlarge":  {HourlyUSD: 1.368, Description: "RDS db.m5.4xlarge MySQL"},
			"db.m6g.large":   {HourlyUSD: 0.152, Description: "RDS db.m6g.large MySQL"},
			"db.m6g.xlarge":  {HourlyUSD: 0.304, Description: "RDS db.m6g.xlarge MySQL"},
			"db.m6g.2xlarge": {HourlyUSD: 0.608, Description: "RDS db.m6g.2xlarge MySQL"},
			"db.r5.large":    {HourlyUSD: 0.240, Description: "RDS db.r5.large MySQL"},
			"db.r5.xlarge":   {HourlyUSD: 0.480, Description: "RDS db.r5.xlarge MySQL"},
			"db.r5.2xlarge":  {HourlyUSD: 0.960, Description: "RDS db.r5.2xlarge MySQL"},
			"db.r6g.large":   {HourlyUSD: 0.218, Description: "RDS db.r6g.large MySQL"},
			"db.r6g.xlarge":  {HourlyUSD: 0.435, Description: "RDS db.r6g.xlarge MySQL"},
			"db.r6g.2xlarge": {HourlyUSD: 0.871, Description: "RDS db.r6g.2xlarge MySQL"},
		},

		// ── Aurora clusters ───────────────────────────────────────────────────
		"aws_rds_cluster": {
			"db.t3.medium":   {HourlyUSD: 0.082, Description: "Aurora db.t3.medium MySQL"},
			"db.t4g.medium":  {HourlyUSD: 0.073, Description: "Aurora db.t4g.medium MySQL"},
			"db.r5.large":    {HourlyUSD: 0.29, Description: "Aurora db.r5.large MySQL"},
			"db.r5.xlarge":   {HourlyUSD: 0.58, Description: "Aurora db.r5.xlarge MySQL"},
			"db.r5.2xlarge":  {HourlyUSD: 1.16, Description: "Aurora db.r5.2xlarge MySQL"},
			"db.r6g.large":   {HourlyUSD: 0.26, Description: "Aurora db.r6g.large MySQL"},
			"db.r6g.xlarge":  {HourlyUSD: 0.52, Description: "Aurora db.r6g.xlarge MySQL"},
			"db.r6g.2xlarge": {HourlyUSD: 1.04, Description: "Aurora db.r6g.2xlarge MySQL"},
			"db.serverless":  {HourlyUSD: 0.12, Description: "Aurora Serverless v2 (per ACU-hr)"},
		},

		// ── ElastiCache ───────────────────────────────────────────────────────
		"aws_elasticache_cluster": {
			"cache.t3.micro":  {HourlyUSD: 0.017, Description: "ElastiCache cache.t3.micro"},
			"cache.t3.small":  {HourlyUSD: 0.034, Description: "ElastiCache cache.t3.small"},
			"cache.t3.medium": {HourlyUSD: 0.068, Description: "ElastiCache cache.t3.medium"},
			"cache.t4g.micro": {HourlyUSD: 0.016, Description: "ElastiCache cache.t4g.micro"},
			"cache.t4g.small": {HourlyUSD: 0.032, Description: "ElastiCache cache.t4g.small"},
			"cache.m5.large":  {HourlyUSD: 0.156, Description: "ElastiCache cache.m5.large"},
			"cache.m5.xlarge": {HourlyUSD: 0.313, Description: "ElastiCache cache.m5.xlarge"},
			"cache.m6g.large":  {HourlyUSD: 0.136, Description: "ElastiCache cache.m6g.large"},
			"cache.m6g.xlarge": {HourlyUSD: 0.272, Description: "ElastiCache cache.m6g.xlarge"},
			"cache.r5.large":  {HourlyUSD: 0.228, Description: "ElastiCache cache.r5.large"},
			"cache.r5.xlarge": {HourlyUSD: 0.455, Description: "ElastiCache cache.r5.xlarge"},
			"cache.r6g.large":  {HourlyUSD: 0.206, Description: "ElastiCache cache.r6g.large"},
			"cache.r6g.xlarge": {HourlyUSD: 0.413, Description: "ElastiCache cache.r6g.xlarge"},
		},
		"aws_elasticache_replication_group": {
			"cache.t3.micro":   {HourlyUSD: 0.017, Description: "ElastiCache cache.t3.micro (repl group)"},
			"cache.t3.small":   {HourlyUSD: 0.034, Description: "ElastiCache cache.t3.small (repl group)"},
			"cache.t3.medium":  {HourlyUSD: 0.068, Description: "ElastiCache cache.t3.medium (repl group)"},
			"cache.m5.large":   {HourlyUSD: 0.156, Description: "ElastiCache cache.m5.large (repl group)"},
			"cache.r5.large":   {HourlyUSD: 0.228, Description: "ElastiCache cache.r5.large (repl group)"},
			"cache.r6g.large":  {HourlyUSD: 0.206, Description: "ElastiCache cache.r6g.large (repl group)"},
			"cache.r6g.xlarge": {HourlyUSD: 0.413, Description: "ElastiCache cache.r6g.xlarge (repl group)"},
		},

		// ── NAT Gateways (flat rate) ──────────────────────────────────────────
		"aws_nat_gateway": {
			"standard": {HourlyUSD: 0.045, Description: "NAT Gateway ($0.045/hr)"},
		},

		// ── Load Balancers (flat rate, excl. LCU) ─────────────────────────────
		"aws_lb": {
			"standard": {HourlyUSD: 0.0225, Description: "ALB/NLB ($0.0225/hr)"},
		},
		"aws_alb": {
			"standard": {HourlyUSD: 0.0225, Description: "ALB ($0.0225/hr)"},
		},

		// ── EBS volumes ───────────────────────────────────────────────────────
		// Priced per GB-month; stored as hourly = monthly / 730.
		"aws_ebs_volume": {
			"gp2":      {HourlyUSD: 0.10 / HoursPerMonth, Description: "EBS gp2 (per GB-month)"},
			"gp3":      {HourlyUSD: 0.08 / HoursPerMonth, Description: "EBS gp3 (per GB-month)"},
			"io1":      {HourlyUSD: 0.125 / HoursPerMonth, Description: "EBS io1 (per GB-month)"},
			"io2":      {HourlyUSD: 0.125 / HoursPerMonth, Description: "EBS io2 (per GB-month)"},
			"st1":      {HourlyUSD: 0.045 / HoursPerMonth, Description: "EBS st1 (per GB-month)"},
			"sc1":      {HourlyUSD: 0.015 / HoursPerMonth, Description: "EBS sc1 (per GB-month)"},
			"standard": {HourlyUSD: 0.05 / HoursPerMonth, Description: "EBS magnetic (per GB-month)"},
		},

		// ── Elastic IPs ───────────────────────────────────────────────────────
		"aws_eip": {
			"standard": {HourlyUSD: 0.005, Description: "Elastic IP ($0.005/hr when idle)"},
		},

		// ── VPC Endpoints ─────────────────────────────────────────────────────
		"aws_vpc_endpoint": {
			"standard": {HourlyUSD: 0.01, Description: "VPC endpoint ($0.01/hr per AZ)"},
		},

		// ── Redshift ──────────────────────────────────────────────────────────
		"aws_redshift_cluster": {
			"dc2.large":   {HourlyUSD: 0.25, Description: "Redshift dc2.large"},
			"dc2.8xlarge": {HourlyUSD: 4.80, Description: "Redshift dc2.8xlarge"},
			"ra3.xlplus":  {HourlyUSD: 1.086, Description: "Redshift ra3.xlplus"},
			"ra3.4xlarge": {HourlyUSD: 3.26, Description: "Redshift ra3.4xlarge"},
			"ra3.16xlarge": {HourlyUSD: 13.04, Description: "Redshift ra3.16xlarge"},
		},

		// ── OpenSearch / Elasticsearch ────────────────────────────────────────
		"aws_opensearch_domain": {
			"t3.small.search":  {HourlyUSD: 0.036, Description: "OpenSearch t3.small.search"},
			"t3.medium.search": {HourlyUSD: 0.073, Description: "OpenSearch t3.medium.search"},
			"m5.large.search":  {HourlyUSD: 0.167, Description: "OpenSearch m5.large.search"},
			"m5.xlarge.search": {HourlyUSD: 0.334, Description: "OpenSearch m5.xlarge.search"},
			"r5.large.search":  {HourlyUSD: 0.186, Description: "OpenSearch r5.large.search"},
			"r5.xlarge.search": {HourlyUSD: 0.371, Description: "OpenSearch r5.xlarge.search"},
			"m6g.large.search": {HourlyUSD: 0.148, Description: "OpenSearch m6g.large.search"},
		},
		"aws_elasticsearch_domain": {
			"t3.small.elasticsearch":  {HourlyUSD: 0.036, Description: "ES t3.small"},
			"t3.medium.elasticsearch": {HourlyUSD: 0.073, Description: "ES t3.medium"},
			"m5.large.elasticsearch":  {HourlyUSD: 0.167, Description: "ES m5.large"},
			"r5.large.elasticsearch":  {HourlyUSD: 0.186, Description: "ES r5.large"},
		},

		// ── MSK (Kafka) ───────────────────────────────────────────────────────
		"aws_msk_cluster": {
			"kafka.t3.small":  {HourlyUSD: 0.075, Description: "MSK kafka.t3.small"},
			"kafka.m5.large":  {HourlyUSD: 0.21, Description: "MSK kafka.m5.large"},
			"kafka.m5.xlarge": {HourlyUSD: 0.42, Description: "MSK kafka.m5.xlarge"},
			"kafka.m5.2xlarge": {HourlyUSD: 0.84, Description: "MSK kafka.m5.2xlarge"},
		},

		// ── DynamoDB (on-demand, per RCU/WCU) ─────────────────────────────────
		"aws_dynamodb_table": {
			"PAY_PER_REQUEST": {HourlyUSD: 0.0, Description: "DynamoDB on-demand (billed per request, not hourly)"},
			"PROVISIONED":     {HourlyUSD: 0.0065, Description: "DynamoDB provisioned (per WCU-hr approx.)"},
		},

		// ── Lambda ────────────────────────────────────────────────────────────
		"aws_lambda_function": {
			"128":  {HourlyUSD: 0.000000208 * 3600, Description: "Lambda 128 MB (per hr at continuous invocation)"},
			"256":  {HourlyUSD: 0.000000417 * 3600, Description: "Lambda 256 MB (per hr at continuous invocation)"},
			"512":  {HourlyUSD: 0.000000834 * 3600, Description: "Lambda 512 MB (per hr at continuous invocation)"},
			"1024": {HourlyUSD: 0.000001667 * 3600, Description: "Lambda 1024 MB (per hr at continuous invocation)"},
			"2048": {HourlyUSD: 0.000003334 * 3600, Description: "Lambda 2048 MB (per hr at continuous invocation)"},
		},

		// ── ECS ───────────────────────────────────────────────────────────────
		// Fargate pricing: vCPU-hr + GB-hr. Represented as "vCPU-GB" pair.
		"aws_ecs_service": {
			"FARGATE": {HourlyUSD: 0.04048 + 0.004445, Description: "ECS Fargate (per 1 vCPU + 1 GB task)"},
			"EC2":     {HourlyUSD: 0.0, Description: "ECS EC2 launch type (cost is on the EC2 instance)"},
		},

		// ── CloudFront ────────────────────────────────────────────────────────
		"aws_cloudfront_distribution": {
			"standard": {HourlyUSD: 0.0, Description: "CloudFront (billed per request + data transfer)"},
		},

		// ── S3 ────────────────────────────────────────────────────────────────
		"aws_s3_bucket": {
			"standard": {HourlyUSD: 0.023 / HoursPerMonth, Description: "S3 Standard (per GB-month)"},
		},

		// ── SQS / SNS (request-priced, ~$0 hourly) ───────────────────────────
		"aws_sqs_queue": {
			"standard": {HourlyUSD: 0.0, Description: "SQS ($0.40 per 1M requests — not hourly)"},
		},
		"aws_sns_topic": {
			"standard": {HourlyUSD: 0.0, Description: "SNS ($0.50 per 1M publishes — not hourly)"},
		},

		// ── Kinesis ───────────────────────────────────────────────────────────
		"aws_kinesis_stream": {
			"ON_DEMAND":  {HourlyUSD: 0.04, Description: "Kinesis on-demand (per shard-hr)"},
			"PROVISIONED": {HourlyUSD: 0.015, Description: "Kinesis provisioned (per shard-hr)"},
		},

		// ── ECR ───────────────────────────────────────────────────────────────
		"aws_ecr_repository": {
			"standard": {HourlyUSD: 0.10 / HoursPerMonth, Description: "ECR ($0.10 per GB-month)"},
		},

		// ── Transit Gateway ───────────────────────────────────────────────────
		"aws_ec2_transit_gateway": {
			"standard": {HourlyUSD: 0.05, Description: "Transit Gateway ($0.05/hr per attachment)"},
		},
		"aws_ec2_transit_gateway_vpc_attachment": {
			"standard": {HourlyUSD: 0.05, Description: "TGW VPC attachment ($0.05/hr)"},
		},

		// =====================================================================
		//  A Z U R E
		// =====================================================================

		// ── Virtual Machines (on-demand Linux, East US) ───────────────────────
		"azurerm_linux_virtual_machine": {
			// B-series burstable
			"Standard_B1s":  {HourlyUSD: 0.0104, Description: "Azure B1s Linux (East US)"},
			"Standard_B1ms": {HourlyUSD: 0.0207, Description: "Azure B1ms Linux (East US)"},
			"Standard_B2s":  {HourlyUSD: 0.0416, Description: "Azure B2s Linux (East US)"},
			"Standard_B2ms": {HourlyUSD: 0.0832, Description: "Azure B2ms Linux (East US)"},
			"Standard_B4ms": {HourlyUSD: 0.166, Description: "Azure B4ms Linux (East US)"},
			// D-series general-purpose
			"Standard_D2s_v5":  {HourlyUSD: 0.096, Description: "Azure D2s_v5 Linux (East US)"},
			"Standard_D4s_v5":  {HourlyUSD: 0.192, Description: "Azure D4s_v5 Linux (East US)"},
			"Standard_D8s_v5":  {HourlyUSD: 0.384, Description: "Azure D8s_v5 Linux (East US)"},
			"Standard_D16s_v5": {HourlyUSD: 0.768, Description: "Azure D16s_v5 Linux (East US)"},
			"Standard_D32s_v5": {HourlyUSD: 1.536, Description: "Azure D32s_v5 Linux (East US)"},
			"Standard_D2as_v5": {HourlyUSD: 0.086, Description: "Azure D2as_v5 Linux (East US)"},
			"Standard_D4as_v5": {HourlyUSD: 0.173, Description: "Azure D4as_v5 Linux (East US)"},
			"Standard_D8as_v5": {HourlyUSD: 0.346, Description: "Azure D8as_v5 Linux (East US)"},
			// E-series memory-optimised
			"Standard_E2s_v5":  {HourlyUSD: 0.126, Description: "Azure E2s_v5 Linux (East US)"},
			"Standard_E4s_v5":  {HourlyUSD: 0.252, Description: "Azure E4s_v5 Linux (East US)"},
			"Standard_E8s_v5":  {HourlyUSD: 0.504, Description: "Azure E8s_v5 Linux (East US)"},
			"Standard_E16s_v5": {HourlyUSD: 1.008, Description: "Azure E16s_v5 Linux (East US)"},
			"Standard_E2as_v5": {HourlyUSD: 0.113, Description: "Azure E2as_v5 Linux (East US)"},
			"Standard_E4as_v5": {HourlyUSD: 0.226, Description: "Azure E4as_v5 Linux (East US)"},
			// F-series compute-optimised
			"Standard_F2s_v2":  {HourlyUSD: 0.085, Description: "Azure F2s_v2 Linux (East US)"},
			"Standard_F4s_v2":  {HourlyUSD: 0.170, Description: "Azure F4s_v2 Linux (East US)"},
			"Standard_F8s_v2":  {HourlyUSD: 0.340, Description: "Azure F8s_v2 Linux (East US)"},
			"Standard_F16s_v2": {HourlyUSD: 0.680, Description: "Azure F16s_v2 Linux (East US)"},
			// L-series storage-optimised
			"Standard_L8s_v3":  {HourlyUSD: 0.624, Description: "Azure L8s_v3 Linux (East US)"},
			"Standard_L16s_v3": {HourlyUSD: 1.248, Description: "Azure L16s_v3 Linux (East US)"},
			// GPU — NC series
			"Standard_NC6s_v3":   {HourlyUSD: 3.06, Description: "Azure NC6s_v3 (1xV100 GPU, East US)"},
			"Standard_NC12s_v3":  {HourlyUSD: 6.12, Description: "Azure NC12s_v3 (2xV100 GPU, East US)"},
			"Standard_NC24s_v3":  {HourlyUSD: 12.24, Description: "Azure NC24s_v3 (4xV100 GPU, East US)"},
		},
		"azurerm_windows_virtual_machine": {
			"Standard_B1s":     {HourlyUSD: 0.0166, Description: "Azure B1s Windows (East US)"},
			"Standard_B2s":     {HourlyUSD: 0.0624, Description: "Azure B2s Windows (East US)"},
			"Standard_B2ms":    {HourlyUSD: 0.1248, Description: "Azure B2ms Windows (East US)"},
			"Standard_D2s_v5":  {HourlyUSD: 0.136, Description: "Azure D2s_v5 Windows (East US)"},
			"Standard_D4s_v5":  {HourlyUSD: 0.272, Description: "Azure D4s_v5 Windows (East US)"},
			"Standard_D8s_v5":  {HourlyUSD: 0.544, Description: "Azure D8s_v5 Windows (East US)"},
			"Standard_D16s_v5": {HourlyUSD: 1.088, Description: "Azure D16s_v5 Windows (East US)"},
			"Standard_F2s_v2":  {HourlyUSD: 0.119, Description: "Azure F2s_v2 Windows (East US)"},
			"Standard_F4s_v2":  {HourlyUSD: 0.238, Description: "Azure F4s_v2 Windows (East US)"},
			"Standard_F8s_v2":  {HourlyUSD: 0.476, Description: "Azure F8s_v2 Windows (East US)"},
		},
		// Legacy azurerm_virtual_machine maps to Linux prices by default
		"azurerm_virtual_machine": {
			"Standard_B1s":     {HourlyUSD: 0.0104, Description: "Azure B1s Linux (East US)"},
			"Standard_B2s":     {HourlyUSD: 0.0416, Description: "Azure B2s Linux (East US)"},
			"Standard_D2s_v5":  {HourlyUSD: 0.096, Description: "Azure D2s_v5 Linux (East US)"},
			"Standard_D4s_v5":  {HourlyUSD: 0.192, Description: "Azure D4s_v5 Linux (East US)"},
			"Standard_D8s_v5":  {HourlyUSD: 0.384, Description: "Azure D8s_v5 Linux (East US)"},
			"Standard_E2s_v5":  {HourlyUSD: 0.126, Description: "Azure E2s_v5 Linux (East US)"},
			"Standard_F2s_v2":  {HourlyUSD: 0.085, Description: "Azure F2s_v2 Linux (East US)"},
			"Standard_F4s_v2":  {HourlyUSD: 0.170, Description: "Azure F4s_v2 Linux (East US)"},
		},

		// ── AKS cluster (free tier control plane) ─────────────────────────────
		"azurerm_kubernetes_cluster": {
			"Free":     {HourlyUSD: 0.0, Description: "AKS Free tier (control plane)"},
			"Standard": {HourlyUSD: 0.10, Description: "AKS Standard tier ($0.10/hr)"},
			"Premium":  {HourlyUSD: 0.60, Description: "AKS Premium tier ($0.60/hr)"},
		},
		"azurerm_kubernetes_cluster_node_pool": {
			"Standard_D2s_v5":  {HourlyUSD: 0.096, Description: "AKS node D2s_v5"},
			"Standard_D4s_v5":  {HourlyUSD: 0.192, Description: "AKS node D4s_v5"},
			"Standard_D8s_v5":  {HourlyUSD: 0.384, Description: "AKS node D8s_v5"},
			"Standard_D2as_v5": {HourlyUSD: 0.086, Description: "AKS node D2as_v5"},
			"Standard_D4as_v5": {HourlyUSD: 0.173, Description: "AKS node D4as_v5"},
			"Standard_E4s_v5":  {HourlyUSD: 0.252, Description: "AKS node E4s_v5"},
			"Standard_F2s_v2":  {HourlyUSD: 0.085, Description: "AKS node F2s_v2"},
			"Standard_F4s_v2":  {HourlyUSD: 0.170, Description: "AKS node F4s_v2"},
		},

		// ── Azure SQL Database ────────────────────────────────────────────────
		"azurerm_mssql_database": {
			"S0":    {HourlyUSD: 0.0202, Description: "Azure SQL S0 (10 DTU)"},
			"S1":    {HourlyUSD: 0.0403, Description: "Azure SQL S1 (20 DTU)"},
			"S2":    {HourlyUSD: 0.0806, Description: "Azure SQL S2 (50 DTU)"},
			"S3":    {HourlyUSD: 0.1612, Description: "Azure SQL S3 (100 DTU)"},
			"P1":    {HourlyUSD: 0.625, Description: "Azure SQL P1 (125 DTU)"},
			"P2":    {HourlyUSD: 1.25, Description: "Azure SQL P2 (250 DTU)"},
			"GP_S_Gen5_1": {HourlyUSD: 0.1096, Description: "Azure SQL GP Serverless 1 vCore"},
			"GP_S_Gen5_2": {HourlyUSD: 0.2192, Description: "Azure SQL GP Serverless 2 vCores"},
			"GP_Gen5_2":   {HourlyUSD: 0.2877, Description: "Azure SQL GP Provisioned 2 vCores"},
			"GP_Gen5_4":   {HourlyUSD: 0.5754, Description: "Azure SQL GP Provisioned 4 vCores"},
		},

		// ── PostgreSQL Flexible Server ────────────────────────────────────────
		"azurerm_postgresql_flexible_server": {
			"B_Standard_B1ms": {HourlyUSD: 0.0207, Description: "Azure PostgreSQL B1ms"},
			"B_Standard_B2s":  {HourlyUSD: 0.0414, Description: "Azure PostgreSQL B2s"},
			"GP_Standard_D2s_v3": {HourlyUSD: 0.1260, Description: "Azure PostgreSQL D2s_v3"},
			"GP_Standard_D4s_v3": {HourlyUSD: 0.2520, Description: "Azure PostgreSQL D4s_v3"},
			"GP_Standard_D8s_v3": {HourlyUSD: 0.5040, Description: "Azure PostgreSQL D8s_v3"},
			"MO_Standard_E2s_v3": {HourlyUSD: 0.1746, Description: "Azure PostgreSQL E2s_v3"},
			"MO_Standard_E4s_v3": {HourlyUSD: 0.3492, Description: "Azure PostgreSQL E4s_v3"},
		},

		// ── MySQL Flexible Server ─────────────────────────────────────────────
		"azurerm_mysql_flexible_server": {
			"B_Standard_B1ms":    {HourlyUSD: 0.0207, Description: "Azure MySQL B1ms"},
			"B_Standard_B2s":     {HourlyUSD: 0.0414, Description: "Azure MySQL B2s"},
			"GP_Standard_D2ds_v4": {HourlyUSD: 0.126, Description: "Azure MySQL D2ds_v4"},
			"GP_Standard_D4ds_v4": {HourlyUSD: 0.252, Description: "Azure MySQL D4ds_v4"},
			"MO_Standard_E2ds_v4": {HourlyUSD: 0.174, Description: "Azure MySQL E2ds_v4"},
		},

		// ── CosmosDB ──────────────────────────────────────────────────────────
		"azurerm_cosmosdb_account": {
			"standard": {HourlyUSD: 0.008, Description: "CosmosDB (per 100 RU/s-hr, provisioned)"},
		},

		// ── Azure Redis Cache ─────────────────────────────────────────────────
		"azurerm_redis_cache": {
			"C0": {HourlyUSD: 0.022, Description: "Azure Redis C0 (250 MB)"},
			"C1": {HourlyUSD: 0.055, Description: "Azure Redis C1 (1 GB)"},
			"C2": {HourlyUSD: 0.110, Description: "Azure Redis C2 (2.5 GB)"},
			"C3": {HourlyUSD: 0.220, Description: "Azure Redis C3 (6 GB)"},
			"C4": {HourlyUSD: 0.440, Description: "Azure Redis C4 (13 GB)"},
			"C5": {HourlyUSD: 0.881, Description: "Azure Redis C5 (26 GB)"},
			"P1": {HourlyUSD: 0.348, Description: "Azure Redis P1 Premium (6 GB)"},
			"P2": {HourlyUSD: 0.697, Description: "Azure Redis P2 Premium (13 GB)"},
			"P3": {HourlyUSD: 1.394, Description: "Azure Redis P3 Premium (26 GB)"},
		},

		// ── Managed Disks ─────────────────────────────────────────────────────
		"azurerm_managed_disk": {
			"Premium_LRS":  {HourlyUSD: 0.132 / HoursPerMonth, Description: "Azure P10 Premium SSD 128 GiB/mo"},
			"StandardSSD_LRS": {HourlyUSD: 0.075 / HoursPerMonth, Description: "Azure E10 Standard SSD 128 GiB/mo"},
			"Standard_LRS": {HourlyUSD: 0.040 / HoursPerMonth, Description: "Azure S10 Standard HDD 128 GiB/mo"},
		},

		// ── Public IPs ────────────────────────────────────────────────────────
		"azurerm_public_ip": {
			"standard": {HourlyUSD: 0.005, Description: "Azure Static Public IP ($0.005/hr)"},
		},

		// ── App Service Plans ─────────────────────────────────────────────────
		"azurerm_service_plan": {
			"B1":  {HourlyUSD: 0.018, Description: "Azure App Service B1"},
			"B2":  {HourlyUSD: 0.036, Description: "Azure App Service B2"},
			"B3":  {HourlyUSD: 0.071, Description: "Azure App Service B3"},
			"S1":  {HourlyUSD: 0.10, Description: "Azure App Service S1"},
			"S2":  {HourlyUSD: 0.20, Description: "Azure App Service S2"},
			"S3":  {HourlyUSD: 0.40, Description: "Azure App Service S3"},
			"P1v3": {HourlyUSD: 0.124, Description: "Azure App Service P1v3"},
			"P2v3": {HourlyUSD: 0.248, Description: "Azure App Service P2v3"},
			"P3v3": {HourlyUSD: 0.496, Description: "Azure App Service P3v3"},
		},

		// ── Azure Load Balancer ───────────────────────────────────────────────
		"azurerm_lb": {
			"standard": {HourlyUSD: 0.025, Description: "Azure Load Balancer ($0.025/hr)"},
		},

		// ── Azure Application Gateway ─────────────────────────────────────────
		"azurerm_application_gateway": {
			"Standard_v2": {HourlyUSD: 0.246, Description: "Azure App Gateway v2 ($0.246/hr)"},
			"WAF_v2":      {HourlyUSD: 0.443, Description: "Azure WAF App Gateway v2 ($0.443/hr)"},
		},

		// ── Azure NAT Gateway ─────────────────────────────────────────────────
		"azurerm_nat_gateway": {
			"standard": {HourlyUSD: 0.045, Description: "Azure NAT Gateway ($0.045/hr)"},
		},

		// ── Azure Container Registry ──────────────────────────────────────────
		"azurerm_container_registry": {
			"Basic":    {HourlyUSD: 5.0 / HoursPerMonth, Description: "ACR Basic ($5/mo)"},
			"Standard": {HourlyUSD: 20.0 / HoursPerMonth, Description: "ACR Standard ($20/mo)"},
			"Premium":  {HourlyUSD: 50.0 / HoursPerMonth, Description: "ACR Premium ($50/mo)"},
		},

		// ── Azure Storage Account (per-GB, request-based) ─────────────────────
		"azurerm_storage_account": {
			"standard": {HourlyUSD: 0.018 / HoursPerMonth, Description: "Azure Storage Hot LRS (per GB-month)"},
		},

		// ── Azure VPN Gateway ─────────────────────────────────────────────────
		"azurerm_virtual_network_gateway": {
			"VpnGw1":  {HourlyUSD: 0.19, Description: "Azure VPN Gateway VpnGw1"},
			"VpnGw2":  {HourlyUSD: 0.49, Description: "Azure VPN Gateway VpnGw2"},
			"VpnGw3":  {HourlyUSD: 1.25, Description: "Azure VPN Gateway VpnGw3"},
			"ErGw1AZ": {HourlyUSD: 0.361, Description: "Azure ExpressRoute Gateway ErGw1AZ"},
			"ErGw2AZ": {HourlyUSD: 0.907, Description: "Azure ExpressRoute Gateway ErGw2AZ"},
		},

		// ── Azure Firewall ────────────────────────────────────────────────────
		"azurerm_firewall": {
			"standard": {HourlyUSD: 1.25, Description: "Azure Firewall Standard ($1.25/hr)"},
		},

		// ── Azure Service Bus ─────────────────────────────────────────────────
		"azurerm_servicebus_namespace": {
			"Basic":    {HourlyUSD: 0.05 / HoursPerMonth, Description: "Service Bus Basic"},
			"Standard": {HourlyUSD: 9.81 / HoursPerMonth, Description: "Service Bus Standard ($9.81/mo)"},
			"Premium":  {HourlyUSD: 0.928, Description: "Service Bus Premium (per MU-hr)"},
		},

		// =====================================================================
		//  G C P
		// =====================================================================

		// ── Compute Engine (on-demand, us-central1) ───────────────────────────
		"google_compute_instance": {
			// E2 — cost-optimised
			"e2-micro":    {HourlyUSD: 0.00838, Description: "GCE e2-micro (us-central1)"},
			"e2-small":    {HourlyUSD: 0.01675, Description: "GCE e2-small (us-central1)"},
			"e2-medium":   {HourlyUSD: 0.03351, Description: "GCE e2-medium (us-central1)"},
			"e2-standard-2": {HourlyUSD: 0.06701, Description: "GCE e2-standard-2 (us-central1)"},
			"e2-standard-4": {HourlyUSD: 0.13402, Description: "GCE e2-standard-4 (us-central1)"},
			"e2-standard-8": {HourlyUSD: 0.26805, Description: "GCE e2-standard-8 (us-central1)"},
			// N2 — general-purpose
			"n2-standard-2":  {HourlyUSD: 0.0971, Description: "GCE n2-standard-2 (us-central1)"},
			"n2-standard-4":  {HourlyUSD: 0.1942, Description: "GCE n2-standard-4 (us-central1)"},
			"n2-standard-8":  {HourlyUSD: 0.3884, Description: "GCE n2-standard-8 (us-central1)"},
			"n2-standard-16": {HourlyUSD: 0.7769, Description: "GCE n2-standard-16 (us-central1)"},
			"n2-standard-32": {HourlyUSD: 1.5537, Description: "GCE n2-standard-32 (us-central1)"},
			"n2-standard-48": {HourlyUSD: 2.3306, Description: "GCE n2-standard-48 (us-central1)"},
			"n2-standard-64": {HourlyUSD: 3.1074, Description: "GCE n2-standard-64 (us-central1)"},
			// N2D — AMD general-purpose
			"n2d-standard-2":  {HourlyUSD: 0.0845, Description: "GCE n2d-standard-2 (us-central1)"},
			"n2d-standard-4":  {HourlyUSD: 0.1690, Description: "GCE n2d-standard-4 (us-central1)"},
			"n2d-standard-8":  {HourlyUSD: 0.3380, Description: "GCE n2d-standard-8 (us-central1)"},
			"n2d-standard-16": {HourlyUSD: 0.6759, Description: "GCE n2d-standard-16 (us-central1)"},
			// C2 — compute-optimised
			"c2-standard-4":  {HourlyUSD: 0.2088, Description: "GCE c2-standard-4 (us-central1)"},
			"c2-standard-8":  {HourlyUSD: 0.4176, Description: "GCE c2-standard-8 (us-central1)"},
			"c2-standard-16": {HourlyUSD: 0.8352, Description: "GCE c2-standard-16 (us-central1)"},
			"c2-standard-30": {HourlyUSD: 1.5660, Description: "GCE c2-standard-30 (us-central1)"},
			"c2-standard-60": {HourlyUSD: 3.1321, Description: "GCE c2-standard-60 (us-central1)"},
			// N1 — first gen general-purpose
			"n1-standard-1":  {HourlyUSD: 0.0475, Description: "GCE n1-standard-1 (us-central1)"},
			"n1-standard-2":  {HourlyUSD: 0.0950, Description: "GCE n1-standard-2 (us-central1)"},
			"n1-standard-4":  {HourlyUSD: 0.1900, Description: "GCE n1-standard-4 (us-central1)"},
			"n1-standard-8":  {HourlyUSD: 0.3800, Description: "GCE n1-standard-8 (us-central1)"},
			"n1-standard-16": {HourlyUSD: 0.7600, Description: "GCE n1-standard-16 (us-central1)"},
			// N1 high-mem
			"n1-highmem-2":  {HourlyUSD: 0.1184, Description: "GCE n1-highmem-2 (us-central1)"},
			"n1-highmem-4":  {HourlyUSD: 0.2368, Description: "GCE n1-highmem-4 (us-central1)"},
			"n1-highmem-8":  {HourlyUSD: 0.4736, Description: "GCE n1-highmem-8 (us-central1)"},
			"n1-highmem-16": {HourlyUSD: 0.9472, Description: "GCE n1-highmem-16 (us-central1)"},
			// T2A — Arm (Ampere Altra)
			"t2a-standard-1": {HourlyUSD: 0.0385, Description: "GCE t2a-standard-1 (us-central1)"},
			"t2a-standard-2": {HourlyUSD: 0.0770, Description: "GCE t2a-standard-2 (us-central1)"},
			"t2a-standard-4": {HourlyUSD: 0.1540, Description: "GCE t2a-standard-4 (us-central1)"},
			// GPU
			"a2-highgpu-1g": {HourlyUSD: 3.6731, Description: "GCE a2-highgpu-1g (1xA100, us-central1)"},
			"a2-highgpu-2g": {HourlyUSD: 7.3462, Description: "GCE a2-highgpu-2g (2xA100, us-central1)"},
			"a2-highgpu-4g": {HourlyUSD: 14.6924, Description: "GCE a2-highgpu-4g (4xA100, us-central1)"},
		},

		// ── GKE cluster (flat rate per cluster) ───────────────────────────────
		"google_container_cluster": {
			"standard":  {HourlyUSD: 0.10, Description: "GKE Standard cluster ($0.10/hr)"},
			"autopilot": {HourlyUSD: 0.10, Description: "GKE Autopilot cluster ($0.10/hr)"},
		},
		"google_container_node_pool": {
			"e2-standard-2":  {HourlyUSD: 0.06701, Description: "GKE node e2-standard-2"},
			"e2-standard-4":  {HourlyUSD: 0.13402, Description: "GKE node e2-standard-4"},
			"e2-standard-8":  {HourlyUSD: 0.26805, Description: "GKE node e2-standard-8"},
			"n2-standard-2":  {HourlyUSD: 0.0971, Description: "GKE node n2-standard-2"},
			"n2-standard-4":  {HourlyUSD: 0.1942, Description: "GKE node n2-standard-4"},
			"n2-standard-8":  {HourlyUSD: 0.3884, Description: "GKE node n2-standard-8"},
			"n2-standard-16": {HourlyUSD: 0.7769, Description: "GKE node n2-standard-16"},
			"n1-standard-1":  {HourlyUSD: 0.0475, Description: "GKE node n1-standard-1"},
			"n1-standard-2":  {HourlyUSD: 0.0950, Description: "GKE node n1-standard-2"},
			"n1-standard-4":  {HourlyUSD: 0.1900, Description: "GKE node n1-standard-4"},
			"c2-standard-4":  {HourlyUSD: 0.2088, Description: "GKE node c2-standard-4"},
			"c2-standard-8":  {HourlyUSD: 0.4176, Description: "GKE node c2-standard-8"},
		},

		// ── Cloud SQL ─────────────────────────────────────────────────────────
		"google_sql_database_instance": {
			"db-f1-micro":     {HourlyUSD: 0.0105, Description: "Cloud SQL db-f1-micro"},
			"db-g1-small":     {HourlyUSD: 0.0255, Description: "Cloud SQL db-g1-small"},
			"db-n1-standard-1": {HourlyUSD: 0.0500, Description: "Cloud SQL db-n1-standard-1"},
			"db-n1-standard-2": {HourlyUSD: 0.1000, Description: "Cloud SQL db-n1-standard-2"},
			"db-n1-standard-4": {HourlyUSD: 0.2000, Description: "Cloud SQL db-n1-standard-4"},
			"db-n1-standard-8": {HourlyUSD: 0.4000, Description: "Cloud SQL db-n1-standard-8"},
			"db-n1-standard-16": {HourlyUSD: 0.8000, Description: "Cloud SQL db-n1-standard-16"},
			"db-n1-highmem-2": {HourlyUSD: 0.1245, Description: "Cloud SQL db-n1-highmem-2"},
			"db-n1-highmem-4": {HourlyUSD: 0.2490, Description: "Cloud SQL db-n1-highmem-4"},
			"db-n1-highmem-8": {HourlyUSD: 0.4980, Description: "Cloud SQL db-n1-highmem-8"},
		},

		// ── Cloud Spanner ─────────────────────────────────────────────────────
		"google_spanner_instance": {
			"standard": {HourlyUSD: 0.90, Description: "Cloud Spanner ($0.90 per node-hr)"},
		},

		// ── Memorystore Redis ─────────────────────────────────────────────────
		"google_redis_instance": {
			"standard": {HourlyUSD: 0.049, Description: "Memorystore Redis Standard (per GB-hr)"},
		},

		// ── Cloud NAT ─────────────────────────────────────────────────────────
		"google_compute_router_nat": {
			"standard": {HourlyUSD: 0.044, Description: "Cloud NAT ($0.044/hr per gateway)"},
		},

		// ── Cloud Load Balancing ──────────────────────────────────────────────
		"google_compute_forwarding_rule": {
			"standard": {HourlyUSD: 0.025, Description: "GCP forwarding rule ($0.025/hr)"},
		},
		"google_compute_global_forwarding_rule": {
			"standard": {HourlyUSD: 0.025, Description: "GCP global forwarding rule ($0.025/hr)"},
		},

		// ── Static External IPs ───────────────────────────────────────────────
		"google_compute_address": {
			"standard": {HourlyUSD: 0.004, Description: "GCP static IP (unused $0.01/hr, in-use $0.004/hr)"},
		},
		"google_compute_global_address": {
			"standard": {HourlyUSD: 0.004, Description: "GCP global static IP ($0.004/hr)"},
		},

		// ── Persistent Disks ──────────────────────────────────────────────────
		"google_compute_disk": {
			"pd-standard": {HourlyUSD: 0.04 / HoursPerMonth, Description: "GCP pd-standard (per GB-month)"},
			"pd-balanced":  {HourlyUSD: 0.10 / HoursPerMonth, Description: "GCP pd-balanced (per GB-month)"},
			"pd-ssd":       {HourlyUSD: 0.17 / HoursPerMonth, Description: "GCP pd-ssd (per GB-month)"},
		},

		// ── Cloud Storage ─────────────────────────────────────────────────────
		"google_storage_bucket": {
			"standard": {HourlyUSD: 0.020 / HoursPerMonth, Description: "GCS Standard (per GB-month)"},
		},

		// ── BigQuery ──────────────────────────────────────────────────────────
		"google_bigquery_dataset": {
			"standard": {HourlyUSD: 0.0, Description: "BigQuery (billed per TB scanned — $6.25/TB)"},
		},

		// ── Cloud Run ─────────────────────────────────────────────────────────
		"google_cloud_run_service": {
			"standard": {HourlyUSD: 0.00002400 * 3600, Description: "Cloud Run (per vCPU-second)"},
		},
		"google_cloud_run_v2_service": {
			"standard": {HourlyUSD: 0.00002400 * 3600, Description: "Cloud Run v2 (per vCPU-second)"},
		},

		// ── Cloud Functions ───────────────────────────────────────────────────
		"google_cloudfunctions_function": {
			"128":  {HourlyUSD: 0.000000231 * 3600, Description: "Cloud Function 128 MB (per hr continuous)"},
			"256":  {HourlyUSD: 0.000000463 * 3600, Description: "Cloud Function 256 MB (per hr continuous)"},
			"512":  {HourlyUSD: 0.000000925 * 3600, Description: "Cloud Function 512 MB (per hr continuous)"},
			"1024": {HourlyUSD: 0.000001650 * 3600, Description: "Cloud Function 1024 MB (per hr continuous)"},
			"2048": {HourlyUSD: 0.000003300 * 3600, Description: "Cloud Function 2048 MB (per hr continuous)"},
		},
		"google_cloudfunctions2_function": {
			"128":  {HourlyUSD: 0.000000231 * 3600, Description: "Cloud Function v2 128 MB (per hr continuous)"},
			"256":  {HourlyUSD: 0.000000463 * 3600, Description: "Cloud Function v2 256 MB (per hr continuous)"},
			"512":  {HourlyUSD: 0.000000925 * 3600, Description: "Cloud Function v2 512 MB (per hr continuous)"},
			"1024": {HourlyUSD: 0.000001650 * 3600, Description: "Cloud Function v2 1024 MB (per hr continuous)"},
		},

		// ── Pub/Sub ───────────────────────────────────────────────────────────
		"google_pubsub_topic": {
			"standard": {HourlyUSD: 0.0, Description: "Pub/Sub ($40 per TiB ingested — not hourly)"},
		},

		// ── VPN ───────────────────────────────────────────────────────────────
		"google_compute_vpn_gateway": {
			"standard": {HourlyUSD: 0.075, Description: "GCP Classic VPN ($0.075/hr)"},
		},
		"google_compute_ha_vpn_gateway": {
			"standard": {HourlyUSD: 0.075, Description: "GCP HA VPN ($0.075/hr per tunnel)"},
		},

		// ── GCR / Artifact Registry ───────────────────────────────────────────
		"google_artifact_registry_repository": {
			"standard": {HourlyUSD: 0.10 / HoursPerMonth, Description: "Artifact Registry ($0.10 per GB-month)"},
		},

		// ── Cloud Armor ───────────────────────────────────────────────────────
		"google_compute_security_policy": {
			"standard": {HourlyUSD: 0.005, Description: "Cloud Armor policy ($5/mo + $1/rule/mo)"},
		},
	}
}
