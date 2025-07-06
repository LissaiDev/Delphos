import React from "react";
import { Host } from "@/types/monitor";
import { formatUptime } from "@/utils/formatters";
import Icon from "@/components/ui/Icon";
import { motion } from "framer-motion";

interface SystemOverviewProps {
  host: Host;
}

export default React.memo(function SystemOverview({
  host,
}: SystemOverviewProps) {
  const overviewItems = [
    {
      label: "Hostname",
      value: host.hostname,
      color: "purple",
      icon: "computer",
    },
    {
      label: "Operating System",
      value: host.os,
      color: "blue",
      icon: "cpu",
    },
    {
      label: "Uptime",
      value: formatUptime(host.uptime),
      color: "green",
      icon: "clock",
    },
  ];

  return (
    <motion.div
      className="bg-slate-800/50 rounded-xl p-6 mb-8 border border-slate-700/50"
      initial={{ opacity: 0, y: 20, scale: 0.95 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{ duration: 0.4, ease: "easeOut" }}
    >
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {overviewItems.map((item) => (
          <div key={item.label} className="flex items-center space-x-4">
            <div className={`p-3 bg-${item.color}-600/20 rounded-lg`}>
              <Icon
                name={item.icon}
                className={`text-${item.color}-400`}
                size="lg"
              />
            </div>
            <div>
              <p className="text-slate-400 text-sm">{item.label}</p>
              <p className="text-white font-semibold text-sm">{item.value}</p>
            </div>
          </div>
        ))}
      </div>
    </motion.div>
  );
});
