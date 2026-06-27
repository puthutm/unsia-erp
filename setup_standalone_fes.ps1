# Setup Standalone Frontend Repositories for UNSIA ERP Modules
# This script copies base configs, components, contexts, and libraries from unsia-portal-web
# and maps the CRM, HRIS, Assessment, Reference, and Academic pages to standalone routes.

$RootPath = "d:\Superman\Superman\Coding\New folder\Dockument\candi\unsia-docs-md"
$SourcePath = "$RootPath\frontend\unsia-portal-web"

$Modules = @(
    @{
        Name = "unsia-crm"
        Pages = @(
            @{ Src = "crm/page.tsx"; Dest = "app/page.tsx" }
            @{ Src = "crm/leads/page.tsx"; Dest = "app/leads/page.tsx" }
        )
        Hooks = @("use-crm.ts")
    },
    @{
        Name = "unsia-hris"
        Pages = @(
            @{ Src = "hris/page.tsx"; Dest = "app/page.tsx" }
            @{ Src = "hris/leave/page.tsx"; Dest = "app/leave/page.tsx" }
        )
        Hooks = @("use-hris.ts")
    },
    @{
        Name = "unsia-assessment"
        Pages = @(
            @{ Src = "assessment/page.tsx"; Dest = "app/page.tsx" }
            @{ Src = "assessment/exam/page.tsx"; Dest = "app/exam/page.tsx" }
        )
        Hooks = @("use-assessment.ts")
    },
    @{
        Name = "unsia-reference"
        Pages = @(
            @{ Src = "reference/page.tsx"; Dest = "app/page.tsx" }
        )
        Hooks = @()
    },
    @{
        Name = "unsia-academic"
        Pages = @(
            @{ Src = "academic/page.tsx"; Dest = "app/page.tsx" }
            @{ Src = "academic/student/page.tsx"; Dest = "app/student/page.tsx" }
            @{ Src = "academic/krs/page.tsx"; Dest = "app/krs/page.tsx" }
            @{ Src = "academic/schedule/page.tsx"; Dest = "app/schedule/page.tsx" }
            @{ Src = "academic/grade/page.tsx"; Dest = "app/grade/page.tsx" }
            @{ Src = "academic/transcript/page.tsx"; Dest = "app/transcript/page.tsx" }
        )
        Hooks = @("use-academic.ts")
    }
)

# Shared configs to copy
$ConfigFiles = @(
    "package.json",
    "tsconfig.json",
    "tailwind.config.js",
    "postcss.config.js",
    "next.config.js",
    "next-env.d.ts",
    "middleware.ts"
)

foreach ($Mod in $Modules) {
    $DestModPath = "$RootPath\frontend\$($Mod.Name)"
    Write-Host "Setting up Standalone Repo: $($Mod.Name) at $DestModPath"

    # Ensure directories exist
    New-Item -ItemType Directory -Force -Path $DestModPath | Out-Null
    New-Item -ItemType Directory -Force -Path "$DestModPath\app" | Out-Null
    New-Item -ItemType Directory -Force -Path "$DestModPath\components" | Out-Null
    New-Item -ItemType Directory -Force -Path "$DestModPath\contexts" | Out-Null
    New-Item -ItemType Directory -Force -Path "$DestModPath\hooks" | Out-Null
    New-Item -ItemType Directory -Force -Path "$DestModPath\lib" | Out-Null

    # Copy configurations
    foreach ($File in $ConfigFiles) {
        if (Test-Path "$SourcePath\$File") {
            Copy-Item -Path "$SourcePath\$File" -Destination "$DestModPath\$File" -Force
        }
    }

    # Copy shared folders
    if (Test-Path "$SourcePath\components") {
        Copy-Item -Path "$SourcePath\components\*" -Destination "$DestModPath\components" -Recurse -Force
    }
    if (Test-Path "$SourcePath\contexts") {
        Copy-Item -Path "$SourcePath\contexts\*" -Destination "$DestModPath\contexts" -Recurse -Force
    }
    if (Test-Path "$SourcePath\lib") {
        Copy-Item -Path "$SourcePath\lib\*" -Destination "$DestModPath\lib" -Recurse -Force
    }
    if (Test-Path "$SourcePath\app\globals.css") {
        Copy-Item -Path "$SourcePath\app\globals.css" -Destination "$DestModPath\app\globals.css" -Force
    }

    # Copy specific pages
    foreach ($Page in $Mod.Pages) {
        $SrcPagePath = "$SourcePath\app\(portal)\$($Page.Src)"
        $DestPagePath = "$DestModPath\$($Page.Dest)"
        
        # Ensure parent directory of page exists
        $DestParentDir = Split-Path -Path $DestPagePath
        New-Item -ItemType Directory -Force -Path $DestParentDir | Out-Null

        if (Test-Path $SrcPagePath) {
            Copy-Item -Path $SrcPagePath -Destination $DestPagePath -Force
            Write-Host "  Copied page: $($Page.Src) -> $($Page.Dest)"
        }
    }

    # Copy specific hooks
    foreach ($Hook in $Mod.Hooks) {
        if (Test-Path "$SourcePath\hooks\$Hook") {
            Copy-Item -Path "$SourcePath\hooks\$Hook" -Destination "$DestModPath\hooks\$Hook" -Force
            Write-Host "  Copied hook: $Hook"
        }
    }

    # Create root layout.tsx
    $LayoutContent = @"
import "./globals.css";
import { AuthProvider } from "@/contexts/auth-context";
import { ReferenceProvider } from "@/contexts/reference-context";

export const metadata = {
  title: "$($Mod.Name.ToUpper()) - UNSIA ERP",
  description: "Standalone Module for $($Mod.Name)",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="id">
      <body>
        <AuthProvider>
          <ReferenceProvider>
            <div className="min-h-screen bg-gray-50 p-6">
              {children}
            </div>
          </ReferenceProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
"@
    [System.IO.File]::WriteAllText("$DestModPath\app\layout.tsx", $LayoutContent)

    # Create hooks/index.ts for the specific repo
    $IndexContent = ""
    if ($Mod.Hooks.Count -gt 0) {
        $IndexContent += "export { useAuth } from '@/contexts/auth-context';`n"
        $IndexContent += "export { useReference } from '@/contexts/reference-context';`n"
        foreach ($Hook in $Mod.Hooks) {
            $HookBase = $Hook.Replace(".ts", "")
            $HookFunc = "use" + ($HookBase.Replace("use-", "").ToUpper())
            # Map use-crm -> useCRM, use-hris -> useHRIS, use-assessment -> useAssessment, use-academic -> useAcademic
            if ($HookBase -eq "use-crm") { $HookFunc = "useCRM" }
            if ($HookBase -eq "use-hris") { $HookFunc = "useHRIS" }
            if ($HookBase -eq "use-assessment") { $HookFunc = "useAssessment" }
            if ($HookBase -eq "use-academic") { $HookFunc = "useAcademic" }
            $IndexContent += "export { $HookFunc } from './$HookBase';`n"
        }
    } else {
        $IndexContent = "export { useAuth } from '@/contexts/auth-context';`nexport { useReference } from '@/contexts/reference-context';`n"
    }
    [System.IO.File]::WriteAllText("$DestModPath\hooks\index.ts", $IndexContent)
}

Write-Host "All standalone frontend repositories set up successfully!"
