import os
import shutil

root_path = r"d:\Superman\Superman\Coding\New folder\Dockument\candi\unsia-docs-md"
source_path = os.path.join(root_path, "frontend", "unsia-portal-web")

modules = [
    {
        "name": "unsia-crm",
        "pages": [
            {"src": "crm/page.tsx", "dest": "app/page.tsx"},
            {"src": "crm/leads/page.tsx", "dest": "app/leads/page.tsx"}
        ],
        "hooks": ["use-crm.ts"]
    },
    {
        "name": "unsia-hris",
        "pages": [
            {"src": "hris/page.tsx", "dest": "app/page.tsx"},
            {"src": "hris/leave/page.tsx", "dest": "app/leave/page.tsx"}
        ],
        "hooks": ["use-hris.ts"]
    },
    {
        "name": "unsia-assessment",
        "pages": [
            {"src": "assessment/page.tsx", "dest": "app/page.tsx"},
            {"src": "assessment/exam/page.tsx", "dest": "app/exam/page.tsx"}
        ],
        "hooks": ["use-assessment.ts"]
    },
    {
        "name": "unsia-reference",
        "pages": [
            {"src": "reference/page.tsx", "dest": "app/page.tsx"}
        ],
        "hooks": []
    },
    {
        "name": "unsia-academic",
        "pages": [
            {"src": "academic/page.tsx", "dest": "app/page.tsx"},
            {"src": "academic/student/page.tsx", "dest": "app/student/page.tsx"},
            {"src": "academic/krs/page.tsx", "dest": "app/krs/page.tsx"},
            {"src": "academic/schedule/page.tsx", "dest": "app/schedule/page.tsx"},
            {"src": "academic/grade/page.tsx", "dest": "app/grade/page.tsx"},
            {"src": "academic/transcript/page.tsx", "dest": "app/transcript/page.tsx"}
        ],
        "hooks": ["use-academic.ts"]
    },
    {
        "name": "unsia-pmb",
        "pages": [],
        "hooks": []
    }
]

config_files = [
    "package.json",
    "tsconfig.json",
    "tailwind.config.js",
    "postcss.config.js",
    "next.config.js",
    "next-env.d.ts",
    "middleware.ts"
]

def setup():
    for mod in modules:
        dest_mod_path = os.path.join(root_path, "frontend", mod["name"])
        print(f"Setting up Standalone Repo: {mod['name']} at {dest_mod_path}")
        
        # Ensure directories
        os.makedirs(os.path.join(dest_mod_path, "app"), exist_ok=True)
        os.makedirs(os.path.join(dest_mod_path, "components"), exist_ok=True)
        os.makedirs(os.path.join(dest_mod_path, "contexts"), exist_ok=True)
        os.makedirs(os.path.join(dest_mod_path, "hooks"), exist_ok=True)
        os.makedirs(os.path.join(dest_mod_path, "lib"), exist_ok=True)
        
        # Copy configurations
        for file in config_files:
            src_file = os.path.join(source_path, file)
            if os.path.exists(src_file):
                shutil.copy2(src_file, os.path.join(dest_mod_path, file))
                
        # Copy shared directories
        for folder in ["components", "contexts", "lib"]:
            src_folder = os.path.join(source_path, folder)
            if os.path.exists(src_folder):
                dest_folder = os.path.join(dest_mod_path, folder)
                if os.path.exists(dest_folder):
                    shutil.rmtree(dest_folder)
                shutil.copytree(src_folder, dest_folder)
                
        # Copy global css
        src_css = os.path.join(source_path, "app", "globals.css")
        if os.path.exists(src_css):
            shutil.copy2(src_css, os.path.join(dest_mod_path, "app", "globals.css"))
            
        # Copy specific pages
        for page in mod["pages"]:
            src_page = os.path.join(source_path, "app", "(portal)", page["src"])
            dest_page = os.path.join(dest_mod_path, page["dest"])
            os.makedirs(os.path.dirname(dest_page), exist_ok=True)
            if os.path.exists(src_page):
                shutil.copy2(src_page, dest_page)
                print(f"  Copied page: {page['src']} -> {page['dest']}")
                
        # Copy specific hooks
        for hook in mod["hooks"]:
            src_hook = os.path.join(source_path, "hooks", hook)
            if os.path.exists(src_hook):
                shutil.copy2(src_hook, os.path.join(dest_mod_path, "hooks", hook))
                print(f"  Copied hook: {hook}")
                
        # Create root layout.tsx
        layout_content = f"""import "./globals.css";
import {{ AuthProvider }} from "@/contexts/auth-context";
import {{ ReferenceProvider }} from "@/contexts/reference-context";

export const metadata = {{
  title: "{mod['name'].upper()} - UNSIA ERP",
  description: "Standalone Module for {mod['name']}",
}};

export default function RootLayout({{
  children,
}}: {{
  children: React.ReactNode;
}}) {{
  return (
    <html lang="id">
      <body>
        <AuthProvider>
          <ReferenceProvider>
            <div className="min-h-screen bg-gray-50 p-6">
              {{children}}
            </div>
          </ReferenceProvider>
        </AuthProvider>
      </body>
    </html>
  );
}}
"""
        with open(os.path.join(dest_mod_path, "app", "layout.tsx"), "w", encoding="utf-8") as f:
            f.write(layout_content)
            
        # Create hooks/index.ts
        index_content = ""
        if len(mod["hooks"]) > 0:
            index_content += "export { useAuth } from '@/contexts/auth-context';\n"
            index_content += "export { useReference } from '@/contexts/reference-context';\n"
            for hook in mod["hooks"]:
                hook_base = hook.replace(".ts", "")
                hook_func = "use" + hook_base.replace("use-", "").upper()
                if hook_base == "use-crm": hook_func = "useCRM"
                elif hook_base == "use-hris": hook_func = "useHRIS"
                elif hook_base == "use-assessment": hook_func = "useAssessment"
                elif hook_base == "use-academic": hook_func = "useAcademic"
                index_content += f"export {{ {hook_func} }} from './{hook_base}';\n"
        else:
            index_content = "export { useAuth } from '@/contexts/auth-context';\nexport { useReference } from '@/contexts/reference-context';\n"
            
        with open(os.path.join(dest_mod_path, "hooks", "index.ts"), "w", encoding="utf-8") as f:
            f.write(index_content)
            
    print("All standalone frontend repositories set up successfully!")

if __name__ == "__main__":
    setup()
