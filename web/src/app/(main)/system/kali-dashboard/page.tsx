"use client";

import * as React from "react";

import { FolderIcon, PlayIcon } from "lucide-react";
import { useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { api } from "@/lib/api";
import type { KaliCategoryStats } from "@/lib/types";
import { cn } from "@/lib/utils";

export default function KaliDashboardPage() {
  const t = useTranslations("skillsPage");
  const [stats, setStats] = React.useState<KaliCategoryStats[]>([]);
  const [loading, setLoading] = React.useState(true);

  React.useEffect(() => {
    void api
      .kaliStats()
      .then((ss) => setStats(ss))
      .catch(() => console.error("Failed to load Kali stats"))
      .finally(() => setLoading(false));
  }, []);

  const totals = React.useMemo(() => {
    const toolCount = stats.reduce((sum, s) => sum + s.tool_count, 0);
    const skillCount = stats.reduce((sum, s) => sum + s.skill_count, 0);
    const callCount = stats.reduce((sum, s) => sum + s.call_count, 0);
    const findingCount = stats.reduce((sum, s) => sum + s.finding_count, 0);
    return { toolCount, skillCount, callCount, findingCount };
  }, [stats]);

  if (loading) {
    return (
      <div className="flex h-full items-center justify-center">
        <p className="text-sm text-muted-foreground">Loading Kali dashboard…</p>
      </div>
    );
  }

  return (
    <div className="flex flex-1 flex-col overflow-hidden">
      <div className="flex items-center gap-3 border-b px-4 py-2.5 lg:px-6">
        <div className="flex flex-col gap-0.5">
          <h1 className="text-sm font-semibold leading-tight">Kali Category Dashboard</h1>
          <p className="text-muted-foreground text-xs">
            Per-category coverage: tools, installed skills, calls and findings.
          </p>
        </div>
      </div>

      <div className="flex flex-1 overflow-hidden p-4">
        {/* Summary cards */}
        <div className="grid gap-4 mb-6">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between py-2">
              <CardTitle className="text-sm">Total Kali Tools</CardTitle>
              <Badge variant="outline" className="text-xs">
                {totals.toolCount}
              </Badge>
            </CardHeader>
            <CardContent>
              <p className="text-xs text-muted-foreground">
                Across all 16 MITRE ATT&CK style categories
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between py-2">
              <CardTitle className="text-sm">Skills Installed</CardTitle>
              <Badge variant="default" className="text-xs">
                {totals.skillCount}
              </Badge>
            </CardHeader>
            <CardContent>
              <p className="text-xs text-muted-foreground">
                {
                  totals.toolCount > 0
                    ? `${Math.round((totals.skillCount / totals.toolCount) * 100)}% coverage`
                    : "0% coverage"
                }
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between py-2">
              <CardTitle className="text-sm">Skill Calls</CardTitle>
              <Badge variant="default" className="text-xs">
                {totals.callCount}
              </Badge>
            </CardHeader>
            <CardContent>
              <p className="text-xs text-muted-foreground">
                Total skill invocations recorded
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between py-2">
              <CardTitle className="text-sm">Findings</CardTitle>
              <Badge variant="default" className="text-xs">
                {totals.findingCount}
              </Badge>
            </CardHeader>
            <CardContent>
              <p className="text-xs text-muted-foreground">
                Findings linked to Kali categories
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Category table */}
        <div className="rounded-md border">
          <CardHeader>
            <CardTitle className="text-sm">Category Coverage</CardTitle>
            <CardDescription>
              Tool count, installed skills, calls and findings per Kali category.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="text-left text-xs text-muted-foreground border-b">
                    <th className="w-1/4 px-3 py-2">Category</th>
                    <th className="w-1/6 px-3 py-2 text-right">Tools</th>
                    <th className="w-1/6 px-3 py-2 text-right">Skills</th>
                    <th className="w-1/6 px-3 py-2 text-right">Calls</th>
                    <th className="w-1/6 px-3 py-2 text-right">Findings</th>
                    <th className="w-1/6 px-3 py-2">Coverage</th>
                  </tr>
                </thead>
                <tbody>
                  {stats.map((s) => {
                    const pct = s.tool_count > 0 ? Math.round((s.skill_count / s.tool_count) * 100) : 0;
                    return (
                      <tr key={s.name} className="border-b hover:bg-muted/50">
                        <td className="px-3 py-2">
                          <div className="flex items-center gap-2">
                            <FolderIcon className="size-4" />
                            <span className="font-mono">{s.name}</span>
                          </div>
                        </td>
                        <td className="px-3 py-2 text-right">{s.tool_count}</td>
                        <td className="px-3 py-2 text-right">{s.skill_count}</td>
                        <td className="px-3 py-2 text-right">{s.call_count}</td>
                        <td className="px-3 py-2 text-right">{s.finding_count}</td>
                        <td className="px-3 py-2">
                          <div className="flex items-center gap-2">
                            <Badge
                              variant={pct === 100 ? "default" : pct > 0 ? "default" : "outline"}
                              className="text-[10px]"
                            >
                              {pct}%
                            </Badge>
                          </div>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </CardContent>
        </div>

        <div className="mt-4 flex items-center gap-2 text-xs text-muted-foreground">
          <PlayIcon className="size-3.5" />
          <span>
            Install missing skills from the <a href="/system/kali" className="hover:underline">
              Kali Tools Directory
            </a>
          </span>
        </div>
      </div>
    </div>
  );
}
