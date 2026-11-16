"use client";

import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { 
  FileCodeIcon, 
  PackageIcon, 
  SettingsIcon, 
  PlayIcon, 
  DatabaseIcon,
  GlobeIcon,
  TagIcon,
  BookOpenIcon,
  UserIcon,
  AlertCircleIcon,
  CheckCircleIcon
} from "lucide-react";
import { useTranslations } from "next-intl";

interface SkillYAMLDisplayProps {
  skill: any;
}

export function SkillYAMLDisplay({ skill }: SkillYAMLDisplayProps) {
  const t = useTranslations("skillsPage");

  if (!skill.has_yaml) {
    return (
      <div className="flex flex-col items-center justify-center h-full text-center p-8">
        <FileCodeIcon className="size-12 text-muted-foreground mb-4" />
        <h3 className="text-lg font-medium">{t("yamlNotFound")}</h3>
        <p className="text-sm text-muted-foreground mt-2 max-w-md">
          {t("yamlNotFoundDesc")}
        </p>
        <div className="mt-4 text-xs text-muted-foreground font-mono">
          skill.yaml
        </div>
      </div>
    );
  }

  const getStatusVariant = (status: string) => {
    switch (status) {
      case "stable": return "default";
      case "beta": return "secondary";
      case "alpha": return "outline";
      case "deprecated": return "destructive";
      case "archived": return "outline";
      default: return "default";
    }
  };

  const getFirstI18n = (obj: any, fallback = "") => {
    if (!obj) return fallback;
    return obj.en || obj.zh || obj.ko || obj.es || obj.ja || Object.values(obj)[0] || fallback;
  };

  return (
    <ScrollArea className="flex-1">
      <div className="space-y-6 p-4">
        {/* Header Info */}
        <div className="flex flex-col gap-2 border-b pb-4">
          <div className="flex flex-wrap items-center gap-2">
            <h3 className="text-lg font-mono font-semibold">{skill.name}</h3>
            {skill.version && (
              <Badge variant="outline" className="text-xs">v{skill.version}</Badge>
            )}
            {skill.category && (
              <Badge variant="default" className="text-xs">{skill.category}</Badge>
            )}
            {skill.status && (
              <Badge variant={getStatusVariant(skill.status)} className="text-xs">
                {skill.status}
              </Badge>
            )}
          </div>
          <div className="flex flex-wrap gap-1">
            {skill.tags?.map((tag: string) => (
              <Badge key={tag} variant="outline" className="text-xs">{tag}</Badge>
            ))}
          </div>
          {skill.description_i18n && (
            <p className="text-sm text-muted-foreground">{getFirstI18n(skill.description_i18n, skill.description)}</p>
          )}
        </div>

        {/* Basic Info Grid */}
        <div className="grid gap-4 md:grid-cols-2">
          <InfoCard icon={<PackageIcon className="size-4" />} title={t("author")} value={skill.author || "—"} />
          <InfoCard icon={<BookOpenIcon className="size-4" />} title={t("license")} value={skill.license || "—"} />
          {skill.compatibility?.min_artex_version && (
            <InfoCard icon={<SettingsIcon className="size-4" />} title={t("minArtexVersion")} value={skill.compatibility.min_artex_version} />
          )}
          {skill.compatibility?.platforms?.length && (
            <InfoCard icon={<GlobeIcon className="size-4" />} title={t("platforms")} value={skill.compatibility.platforms.join(", ")} />
          )}
          {skill.compatibility?.architectures?.length && (
            <InfoCard icon={<DatabaseIcon className="size-4" />} title={t("architectures")} value={skill.compatibility.architectures.join(", ")} />
          )}
        </div>

        {/* Dependencies */}
        {skill.dependencies && (
          <section className="space-y-3">
            <h4 className="font-medium flex items-center gap-2">
              <PackageIcon className="size-4" />
              {t("dependencies")}
            </h4>
            <div className="space-y-2 ml-4">
              {skill.dependencies.system?.length && (
                <DepItem title={t("systemDeps")} items={skill.dependencies.system.map((d: any) => `${d.name}${d.version ? ` >= ${d.version}` : ""}`)} />
              )}
              {skill.dependencies.python?.length && (
                <DepItem title={t("pythonDeps")} items={skill.dependencies.python.map((d: any) => `${d.package}${d.version ? `==${d.version}` : ""}`)} />
              )}
              {skill.dependencies.node?.length && (
                <DepItem title={t("nodeDeps")} items={skill.dependencies.node.map((d: any) => `${d.package}${d.version ? `@${d.version}` : ""}`)} />
              )}
              {skill.dependencies.go?.length && (
                <DepItem title={t("goDeps")} items={skill.dependencies.go.map((d: any) => `${d.module}${d.version ? `@${d.version}` : ""}`)} />
              )}
            </div>
          </section>
        )}

        {/* Allowed Tools */}
        {skill.allowed_tools?.length && (
          <section className="space-y-3">
            <h4 className="font-medium flex items-center gap-2">
              <SettingsIcon className="size-4" />
              {t("allowedTools")}
            </h4>
            <div className="flex flex-wrap gap-1">
              {skill.allowed_tools.map((tool: string) => (
                <Badge key={tool} variant="outline" className="text-xs font-mono">{tool}</Badge>
              ))}
            </div>
          </section>
        )}

        {/* Parameters */}
        {skill.parameters?.length && (
          <section className="space-y-3">
            <h4 className="font-medium flex items-center gap-2">
              <SettingsIcon className="size-4" />
              {t("parameters")}
            </h4>
            <div className="space-y-2 ml-4">
              {skill.parameters.map((p: any, i: number) => (
                <ParamItem key={i} param={p} />
              ))}
            </div>
          </section>
        )}

        {/* Entrypoint */}
        {skill.entrypoint && (
          <section className="space-y-3">
            <h4 className="font-medium flex items-center gap-2">
              <PlayIcon className="size-4" />
              {t("entrypoint")}
            </h4>
            <div className="space-y-2 ml-4 font-mono text-sm">
              <div className="flex items-center gap-2">
                <Label className="text-xs text-muted-foreground w-20">{t("type")}</Label>
                <Badge variant="secondary">{skill.entrypoint.type}</Badge>
              </div>
              <div className="flex items-center gap-2">
                <Label className="text-xs text-muted-foreground w-20">{t("command")}</Label>
                <code>{skill.entrypoint.command}</code>
              </div>
              {skill.entrypoint.args?.length && (
                <div className="flex items-center gap-2">
                  <Label className="text-xs text-muted-foreground w-20">{t("args")}</Label>
                  <code>{skill.entrypoint.args.join(" ")}</code>
                </div>
              )}
              {skill.entrypoint.working_dir && (
                <div className="flex items-center gap-2">
                  <Label className="text-xs text-muted-foreground w-20">{t("workingDir")}</Label>
                  <code>{skill.entrypoint.working_dir}</code>
                </div>
              )}
            </div>
          </section>
        )}

        {/* Outputs */}
        {skill.outputs?.length && (
          <section className="space-y-3">
            <h4 className="font-medium flex items-center gap-2">
              <DatabaseIcon className="size-4" />
              {t("outputs")}
            </h4>
            <div className="space-y-2 ml-4">
              {skill.outputs.map((o: any, i: number) => (
                <div key={i} className="flex items-center gap-2 p-2 rounded border bg-muted/50">
                  <Badge variant="secondary" className="text-xs">{o.type}</Badge>
                  <span className="font-medium">{o.name}</span>
                  {o.description && <span className="text-xs text-muted-foreground ml-auto">{o.description}</span>}
                </div>
              ))}
            </div>
          </section>
        )}

        {/* Examples */}
        {skill.examples?.length && (
          <section className="space-y-3">
            <h4 className="font-medium flex items-center gap-2">
              <BookOpenIcon className="size-4" />
              {t("examples")}
            </h4>
            <div className="space-y-2 ml-4">
              {skill.examples.map((ex: any, i: number) => (
                <div key={i} className="p-3 rounded border bg-muted/50">
                  <div className="flex items-start gap-2">
                    <Badge variant="outline" className="text-xs">{ex.name}</Badge>
                    <div>
                      <p className="font-medium">{ex.name}</p>
                      {ex.description && <p className="text-xs text-muted-foreground">{ex.description}</p>}
                      {ex.parameters && (
                        <pre className="mt-1 text-xs font-mono bg-muted p-2 rounded overflow-auto">
                          {JSON.stringify(ex.parameters, null, 2)}
                        </pre>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* Maintainers */}
        {skill.maintainers?.length && (
          <section className="space-y-3">
            <h4 className="font-medium flex items-center gap-2">
              <UserIcon className="size-4" />
              {t("maintainers")}
            </h4>
            <div className="flex flex-wrap gap-2">
              {skill.maintainers.map((m: any, i: number) => (
                <Badge key={i} variant="outline" className="text-xs flex items-center gap-1">
                  <UserIcon className="size-3" />
                  {m.name}{m.email && ` <${m.email}>`}
                </Badge>
              ))}
            </div>
          </section>
        )}

        {/* Documentation Links */}
        {skill.documentation && (
          <section className="space-y-3">
            <h4 className="font-medium flex items-center gap-2">
              <BookOpenIcon className="size-4" />
              {t("documentation")}
            </h4>
            <div className="flex flex-wrap gap-2">
              {skill.documentation.readme && (
                <a href="#" className="text-sm text-primary hover:underline flex items-center gap-1">
                  <FileCodeIcon className="size-3" />
                  {t("readme")}
                </a>
              )}
              {skill.documentation.wiki && (
                <a href={skill.documentation.wiki} target="_blank" rel="noopener noreferrer" className="text-sm text-primary hover:underline flex items-center gap-1">
                  <GlobeIcon className="size-3" />
                  {t("wiki")}
                </a>
              )}
              {skill.documentation.changelog && (
                <a href="#" className="text-sm text-primary hover:underline flex items-center gap-1">
                  <BookOpenIcon className="size-3" />
                  {t("changelog")}
                </a>
              )}
            </div>
          </section>
        )}
      </div>
    </ScrollArea>
  );
}

function InfoCard({ icon, title, value }: { icon: React.ReactNode; title: string; value: string }) {
  const t = useTranslations("skillsPage");
  return (
    <div className="flex items-start gap-2 p-3 rounded border bg-muted/50">
      <div className="text-muted-foreground mt-0.5">{icon}</div>
      <div className="flex-1 min-w-0">
        <p className="text-xs text-muted-foreground">{t(title.toLowerCase())}</p>
        <p className="font-mono text-sm truncate">{value}</p>
      </div>
    </div>
  );
}

function DepItem({ title, items }: { title: string; items: string[] }) {
  const t = useTranslations("skillsPage");
  return (
    <div className="flex items-start gap-2 p-2 rounded border bg-muted/50">
      <Label className="text-xs text-muted-foreground w-24">{t(title.toLowerCase())}</Label>
      <div className="flex flex-wrap gap-1">
        {items.map((item, i) => (
          <Badge key={i} variant="outline" className="text-xs font-mono">{item}</Badge>
        ))}
      </div>
    </div>
  );
}

function ParamItem({ param }: { param: any }) {
  const t = useTranslations("skillsPage");
  const required = param.required ? <span className="text-destructive ml-1">*</span> : null;
  return (
    <div className="p-3 rounded border bg-muted/50">
      <div className="flex items-center gap-2 mb-1">
        <span className="font-medium font-mono">{param.name}</span>
        {required}
        <Badge variant="secondary" className="text-xs">{param.type}</Badge>
        {param.secret && (
              <span title={t("secretParam")} className="cursor-help">
                <AlertCircleIcon className="size-3 text-amber-500" />
              </span>
            )}
      </div>
      <p className="text-xs text-muted-foreground ml-6">{getFirstI18n(param.description) || getFirstI18n(param.label)}</p>
      {param.enum?.length && (
        <div className="mt-1 ml-6 flex flex-wrap gap-1">
          {param.enum.map((e: any, i: number) => (
            <Badge key={i} variant="outline" className="text-xs">{e.value}{e.label && `: ${getFirstI18n(e.label)}`}</Badge>
          ))}
        </div>
      )}
      {param.validation && (
        <div className="mt-1 ml-6 text-xs text-muted-foreground">
          {param.validation.pattern && <span>Pattern: {param.validation.pattern}</span>}
          {param.validation.min !== undefined && <span className="mx-2">Min: {param.validation.min}</span>}
          {param.validation.max !== undefined && <span className="mx-2">Max: {param.validation.max}</span>}
          {param.validation.min_length !== undefined && <span className="mx-2">MinLen: {param.validation.min_length}</span>}
          {param.validation.max_length !== undefined && <span className="mx-2">MaxLen: {param.validation.max_length}</span>}
        </div>
      )}
    </div>
  );
}

function getFirstI18n(obj: any, fallback = "") {
  if (!obj) return fallback;
  return obj.en || obj.zh || obj.ko || obj.es || obj.ja || Object.values(obj)[0] || fallback;
}