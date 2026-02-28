Implementation Flow Of Centralized OpenTelemetry Infra Integration & Consumption Steps

PART 1 — Centralized Infra Repository
Repository: mhc-infra-observability
Objective: Finalize observability infra, tag a stable release, and make it consumable by services.

STEP 1 :  Navigate to infra repository
cd mhc-infra-observability
STEP 2 : Verify module health
go mod tidy
go build ./…

Must complete with zero errors.
STEP 3: Commit stabilized changes
git add .
git commit -m "chore(observability): stabilize OpenTelemetry init and version alignment"
STEP 4 : Tag and push release (REQUIRED)
git tag v0.1.4
git push origin main
git push origin v0.1.4
Tags are created only in infra repositories, never in service repositories.
5. Optional sanity check
go list -m all | findstr opentelemetry



NOTE:
When ever we do the changes in this Infra observability repo we must follow this above all the steps.
 After that update the tag, in service repo through following steps below for sure.
PART 2 — All Service Repositories [tenant, identity..]
Repository: mhc-backend-lms-tenant-subscription-service
Objective: Consume tagged infra repo and run service with OpenTelemetry enabled.

STEP 6 : Navigate to service repository
cd mhc-backend-lms-tenant-subscription-service
STEP 7 : Configure private module access (one-time per machine)
go env -w GOPRIVATE=github.com/MH-Cognition/*
STEP 8. Update observability dependency (use a single tagged version)
go get github.com/MH-Cognition/mhc-infra-observability@v0.1.6
Always depend on a tagged version, never main.
STEP 9 : Clean and sync dependencies
go clean -modcache
go mod tidy
STEP 10. Build service
go build ./...

 Outcome
Observability infra is versioned and immutable


Tenant service consumes a stable OpenTelemetry setup


Traces, metrics, and logs are consistently initialized
