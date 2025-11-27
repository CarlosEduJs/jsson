"use client";

import { useState, useEffect } from "react";
import { CodeBlock } from "../shared/code-block";
import { Tabs, TabsList, TabsTrigger } from "../ui/tabs";
import { Button } from "../ui/button";
import { Copy } from "lucide-react";
import { toastManager } from "@/components/ui/toast";

interface OutputViewerProps {
  jssonCode: string;
  error?: string | null;
}

export function OutputViewer({ jssonCode, error }: OutputViewerProps) {
  const [format, setFormat] = useState("json");
  const [output, setOutput] = useState("");
  const [transpileError, setTranspileError] = useState<string | null>(null);

  useEffect(() => {
    if (!jssonCode || error) {
      setOutput("");
      return;
    }

    if (window.transpileJSSON) {
      const result = window.transpileJSSON(jssonCode, format);
      if (result.error) {
        setTranspileError(result.error);
        setOutput("");
      } else {
        setTranspileError(null);
        setOutput(result.output || "");
      }
    }
  }, [jssonCode, format, error]);

  const displayError = error || transpileError;

  function copyToClipboard() {
    navigator.clipboard.writeText(output);

    toastManager.add({
      title: "Copied to clipboard",
      description: "The output has been copied to your clipboard.",
    });
  }

  return (
    <div className="h-full w-full flex flex-col rounded-xl overflow-hidden border border-border bg-card/50 backdrop-blur-sm shadow-sm">
      <div className="flex items-center justify-between px-2 py-2 border-b border-border bg-muted/30">
        <Tabs value={format} onValueChange={setFormat} className="w-full">
          <TabsList className="h-8 bg-transparent p-0 gap-2">
            <TabsTrigger
              value="json"
              className="data-[state=active]:bg-background data-[state=active]:shadow-sm h-7 text-xs px-3 rounded-md"
            >
              JSON
            </TabsTrigger>
            <TabsTrigger
              value="yaml"
              className="data-[state=active]:bg-background data-[state=active]:shadow-sm h-7 text-xs px-3 rounded-md"
            >
              YAML
            </TabsTrigger>
            <TabsTrigger
              value="toml"
              className="data-[state=active]:bg-background data-[state=active]:shadow-sm h-7 text-xs px-3 rounded-md"
            >
              TOML
            </TabsTrigger>
            <TabsTrigger
              value="typescript"
              className="data-[state=active]:bg-background data-[state=active]:shadow-sm h-7 text-xs px-3 rounded-md"
            >
              TypeScript
            </TabsTrigger>
          </TabsList>
        </Tabs>
        <Button variant={"ghost"} size={"sm"} onClick={copyToClipboard}>
          <Copy className="h-4 w-4" />
          Copy
        </Button>
      </div>

      <div className="flex-1 overflow-auto bg-card/50 p-0">
        {displayError ? (
          <div className="p-4 text-red-400 font-mono text-sm">
            Error: {displayError}
          </div>
        ) : (
          <CodeBlock
            code={output}
            language={format === "json" ? "json" : "jsson"}
            className="h-full rounded-none border-none bg-transparent"
            showLineNumbers
          />
        )}
      </div>
    </div>
  );
}
