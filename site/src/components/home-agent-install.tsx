import Link from "next/link";
import { CopyCommand } from "@/components/copy-command";
import {
  agentInstallGuide,
  installCommand,
  installInspectCommands,
} from "@/lib/site";

/**
 * Below-fold AI install carrier for URL-drop + "install"
 * (dec-20260808-homepage-ai-install-signals-9eb5a6a9).
 * Secondary to hero — must not own the first viewport.
 */
export function HomeAgentInstall() {
  return (
    <section
      id="for-ai-agents"
      aria-labelledby="for-ai-agents-heading"
      data-pysar-audience="ai-agent"
      data-pysar-intent="install"
      className="size-full overflow-hidden px-4 md:px-5 lg:px-7.5"
    >
      <div className="relative z-10 mx-auto w-full max-w-250 border-x-2 border-zinc-100 py-16 lg:py-24">
        <div className="mx-auto max-w-178 space-y-8 px-5 md:px-12.5">
          <div className="space-y-2.5 text-center md:text-start">
            <div className="inline-flex items-center justify-center rounded-full border-2 border-zinc-200 px-4 py-1.5 text-xs font-semibold uppercase text-primary">
              For AI agents
            </div>
            <h2
              id="for-ai-agents-heading"
              className="text-section md:text-section-md"
            >
              User pasted this URL and said{" "}
              <span className="text-zinc-500">
                &ldquo;{agentInstallGuide.trigger}&rdquo;
              </span>
            </h2>
            <p className="text-zinc-600">
              Follow these steps. Do not invent brew, npm, pip, or{" "}
              <code className="rounded bg-zinc-100 px-1.5 py-0.5 font-mono text-[13px] text-zinc-800">
                go install
              </code>{" "}
              alternatives. Plain-text index:{" "}
              <Link
                href={agentInstallGuide.llmsTxtPath}
                className="font-medium text-zinc-950 underline decoration-primary/50 underline-offset-4 hover:text-primary hover:decoration-primary"
              >
                /llms.txt
              </Link>
              . Human docs:{" "}
              <Link
                href={agentInstallGuide.docsInstallPath}
                className="font-medium text-zinc-950 underline decoration-primary/50 underline-offset-4 hover:text-primary hover:decoration-primary"
              >
                Install
              </Link>{" "}
              and{" "}
              <Link
                href={agentInstallGuide.docsInitPath}
                className="font-medium text-zinc-950 underline decoration-primary/50 underline-offset-4 hover:text-primary hover:decoration-primary"
              >
                Init
              </Link>
              .
            </p>
          </div>

          <ol className="space-y-6 border-y-2 border-zinc-100 py-8 text-sm text-zinc-700">
            <li className="space-y-3">
              <p className="font-semibold text-zinc-900">
                1. Install the{" "}
                <code className="rounded bg-zinc-100 px-1.5 py-0.5 font-mono text-[13px] text-zinc-800">
                  pysar
                </code>{" "}
                binary (macOS / Linux)
              </p>
              <CopyCommand command={installCommand} />
              <p className="text-xs text-zinc-500">
                Or inspect first:{" "}
                <code className="rounded bg-zinc-100 px-1.5 py-0.5 font-mono text-[12px] text-zinc-700">
                  {installInspectCommands.download}
                </code>
                , then{" "}
                <code className="rounded bg-zinc-100 px-1.5 py-0.5 font-mono text-[12px] text-zinc-700">
                  {installInspectCommands.run}
                </code>
                .
              </p>
              <p>
                Confirm with{" "}
                <code className="rounded bg-zinc-100 px-1.5 py-0.5 font-mono text-[13px] text-zinc-800">
                  pysar --version
                </code>
                .
              </p>
            </li>

            <li className="space-y-3">
              <p className="font-semibold text-zinc-900">
                2. Init the writing project for the host you are running in
              </p>
              <ul className="space-y-3">
                {agentInstallGuide.hosts.map((host) => (
                  <li
                    key={host.name}
                    className="rounded border-2 border-zinc-100 bg-zinc-50/80 px-4 py-3"
                  >
                    <p className="mb-2 font-medium text-zinc-900">
                      {host.name}
                    </p>
                    <CopyCommand command={host.initCommand} />
                    <p className="mt-2 text-xs text-zinc-500">{host.note}</p>
                  </li>
                ))}
              </ul>
            </li>

            <li>
              <p className="font-semibold text-zinc-900">
                3. Open the project in that host and continue with{" "}
                <code className="rounded bg-zinc-100 px-1.5 py-0.5 font-mono text-[13px] text-zinc-800">
                  /ps-onboard
                </code>{" "}
                or{" "}
                <code className="rounded bg-zinc-100 px-1.5 py-0.5 font-mono text-[13px] text-zinc-800">
                  /ps
                </code>{" "}
                as the author directs.
              </p>
            </li>
          </ol>
        </div>
      </div>
    </section>
  );
}
