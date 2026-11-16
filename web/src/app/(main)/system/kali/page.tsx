"use client";

import * as React from "react";

import {
  AlertTriangleIcon,
  ChevronRightIcon,
  ExternalLinkIcon,
  FolderIcon,
  PlayIcon,
  PlusIcon,
  Trash2Icon,
  UploadIcon,
} from "lucide-react";
import { useTranslations } from "next-intl";
import { toast } from "sonner";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { api } from "@/lib/api";
import type { KaliCategoryInfo, KaliToolInfo } from "@/lib/types";
import { cn } from "@/lib/utils";

export default function KaliPage() {
  const t = useTranslations("skillsPage");
  const [categories, setCategories] = React.useState<KaliCategoryInfo[]>([]);
  const [tools, setTools] = React.useState<KaliToolInfo[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [selectedCategory, setSelectedCategory] = React.useState<string | null>(null);
  const [installing, setInstalling] = React.useState<string | null>(null);

  React.useEffect(() => {
    void api
      .kaliCategories()
      .then((cs) => setCategories(cs))
      .catch(() => toast.error("Failed to load Kali categories"))
      .finally(() => {
        void api
          .kaliTools()
          .then((ts) => setTools(ts))
          .catch(() => toast.error("Failed to load Kali tools"))
          .finally(() => setLoading(false));
      });
  }, []);

  const filteredTools = selectedCategory
    ? tools.filter((t) => t.category === selectedCategory)
    : tools;

  const toolsByCategory = React.useMemo(() => {
    const map = new Map<string, KaliToolInfo[]>();
    for (const tool of tools) {
      const arr = map.get(tool.category) ?? [];
      arr.push(tool);
      map.set(tool.category, arr);
    }
    return map;
  }, [tools]);

  async function handleInstall(name: string) {
    setInstalling(name);
    try {
      await api.kaliToolInstall(name);
      toast.success(`Installed ${name} as ARTEX skill`);
      // Refresh tools to reflect new skill availability
      const updated = await api.kaliTools();
      setTools(updated);
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e);
      toast.error(`Install failed: ${msg}`);
    } finally {
      setInstalling(null);
    }
  }

  if (loading) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="flex flex-col items-center gap-2">
          <div className="h-8 w-8 rounded-full border-2 border-primary/50 animate-spin"></div>
          <p className="text-sm text-muted-foreground">Loading Kali tools directory…</p>
        </div>
      </div>
    );
  }

  function getInstallButtonContent(toolName: string): React.ReactNode {
    if (installing === toolName) {
      return (
        <span className="flex items-center gap-1">
          <div className="h-3 w-3 rounded-full border-2 border-current animate-spin"></div>
          Installing…
        </span>
      );
    }
    const tool = tools.find((t) => t.name === toolName);
    if (tool?.skill_available) return "Installed";
    return "Install as Skill";
  }

  return (
    <div className="flex flex-1 flex-col overflow-hidden">
      <div className="flex items-center gap-3 border-b px-4 py-2.5 lg:px-6">
        <div className="flex flex-col gap-0.5">
          <h1 className="text-sm font-semibold leading-tight">Kali Linux Tools</h1>
          <p className="text-muted-foreground text-xs">
            Browse Kali Linux tools by category and install them as ARTEX skills.
          </p>
        </div>
      </div>

      <div className="flex flex-1 overflow-hidden">
        {/* 左側：分類選單 */}
        <div className="flex w-72 shrink-0 flex-col border-r">
          <div className="border-b p-4">
            <Label className="text-xs font-medium uppercase tracking-wider">Categories</Label>
            <div className="mt-2 space-y-1">
              <button
                key="all"
                onClick={() => setSelectedCategory(null)}
                className={cn(
                  "flex w-full items-center gap-2 rounded px-2 py-1.5 text-sm",
                  selectedCategory === null ? "bg-accent text-accent-foreground" : "text-muted-foreground hover:bg-muted",
                )}
              >
                <FolderIcon className="size-4" />
                <span className="flex-1 truncate">All ({tools.length})</span>
              </button>
              {categories.map((c) => (
                <button
                  key={c.name}
                  onClick={() => setSelectedCategory(c.name)}
                  className={cn(
                    "flex w-full items-center gap-2 rounded px-2 py-1.5 text-sm",
                    selectedCategory === c.name
                      ? "bg-accent text-accent-foreground"
                      : "text-muted-foreground hover:bg-muted",
                  )}
                >
                  <FolderIcon className="size-4" />
                  <span className="flex-1 truncate">{c.name}</span>
                  <Badge variant="outline" className="text-[10px]">
                    {c.tool_count}
                  </Badge>
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* 右側：工具列表 */}
        <div className="flex flex-1 flex-col overflow-auto p-4">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold">
              {selectedCategory ? categories.find((c) => c.name === selectedCategory)?.name ?? "" : "All Tools"}
            </h2>
            <p className="text-xs text-muted-foreground">
              {filteredTools.length} tool{filteredTools.length !== 1 ? "s" : ""}
            </p>
          </div>

          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {filteredTools.length === 0 ? (
              <div className="col-span-full flex flex-col items-center justify-center py-12 text-muted-foreground">
                <AlertTriangleIcon className="size-12 mb-2" />
                <p>No tools found</p>
              </div>
            ) : (
              filteredTools.map((tool) => (
                <Card key={tool.name} className="hover:border-accent/50 transition-border">
                  <CardHeader className="flex flex-row items-center justify-between py-2">
                    <CardTitle className="text-sm font-mono truncate">{tool.name}</CardTitle>
                    <Badge
                      variant={tool.skill_available ? "default" : "outline"}
                      className="text-[10px] hover:border-accent/50"
                    >
                      {tool.skill_available ? "Installed" : "Kali"}
                    </Badge>
                  </CardHeader>
                  <CardContent className="py-3">
                    <p className="text-xs text-muted-foreground line-clamp-2">{tool.description}</p>
                    <div className="mt-3 flex items-center gap-1">
                      <span className="text-[10px] font-mono text-muted-foreground">
                        {tool.category}
                      </span>
                    </div>
                  </CardContent>
                  <CardFooter>
                    <Button
                      onClick={() => handleInstall(tool.name)}
                      disabled={tool.skill_available || installing === tool.name}
                      className="w-full text-sm"
                      size="sm"
                    >
                      {getInstallButtonContent(tool.name)}
                    </Button>
                    <a
                      href={tool.official_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-xs text-muted-foreground hover:text-primary flex items-center gap-1"
                    >
                      <ExternalLinkIcon className="size-3" />
                      kali.org
                    </a>
                  </CardFooter>
                </Card>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
