import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

const sidebars: SidebarsConfig = {
  docs: [
    "intro",
    {
      type: "category",
      label: "Getting Started",
      items: [
        "getting-started/prerequisites",
        "getting-started/installation",
        "getting-started/configuration",
        "getting-started/running",
      ],
    },
    {
      type: "category",
      label: "Commands",
      items: [
        "commands/text-commands",
        "commands/voice-commands",
      ],
    },
    "architecture",
    "discord-setup",
    {
      type: "category",
      label: "Deployment",
      items: [
        "deployment/docker",
        "deployment/railway",
      ],
    },
  ],
};

export default sidebars;
