import "@/global.css";

// ----------------------------------------------------------------------

import { primaryFont } from "./theme/typography";

import ThemeProvider from "./theme";
import { SettingsProvider } from "@/components/settings";
import { MotionLazy } from "@/components/animate/motion-lazy";
import ProgressBar from "@/components/progress-bar";
import { ApolloWrapper } from "@/apollo-client/apollo-wrapper";
// ----------------------------------------------------------------------

export const viewport = {
  themeColor: "#000000",
  width: "device-width",
  initialScale: 1,
  maximumScale: 1,
};

export const metadata = {
  title: "Checkpoint - dashboard",
  description: "Checkpoint dashboard",
};

type Props = {
  children: React.ReactNode;
};

export default function RootLayout({ children }: Props) {
  return (
    <html lang="en" className={primaryFont.className}>
      <body>
        <ApolloWrapper>
          <SettingsProvider
            defaultSettings={{
              themeMode: "light",
              themeDirection: "ltr",
              themeContrast: "default",
              themeLayout: "vertical",
              themeColorPresets: "blue",
              themeStretch: false,
            }}
          >
            <ThemeProvider>
              <MotionLazy>
                <ProgressBar />
                {children}
              </MotionLazy>
            </ThemeProvider>
          </SettingsProvider>
        </ApolloWrapper>
      </body>
    </html>
  );
}
