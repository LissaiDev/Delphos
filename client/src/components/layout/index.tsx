import { ReactNode } from "react";
import { Titillium_Web, Share_Tech_Mono } from "next/font/google";

const titilliumWeb = Titillium_Web({
  subsets: ["latin"],
  weight: ["200", "300", "400", "600", "700", "900"],
  variable: "--font-titillium-web",
});

const shareTechMono = Share_Tech_Mono({
  subsets: ["latin"],
  weight: ["400"],
  variable: "--font-share-tech-mono",
});

export default function Providers({ children }: { children: ReactNode }) {
  return <div className={`${titilliumWeb.variable} ${shareTechMono.variable} antialiased`}>{children}</div>;
}