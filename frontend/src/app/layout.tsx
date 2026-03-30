import "./globals.css";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Зоомагазин та рослини",
  description: "Sprint 2 frontend"
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="uk">
      <body>{children}</body>
    </html>
  );
}
