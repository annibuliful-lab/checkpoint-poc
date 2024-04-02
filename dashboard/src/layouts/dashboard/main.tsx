import Box, { BoxProps } from "@mui/material/Box";

import { useResponsive } from "@/hooks/use-responsive";
import { HEADER, NAV } from "../config-layout";

// ----------------------------------------------------------------------

const SPACING = 8;

export default function Main({ children, sx, ...other }: BoxProps) {
  const lgUp = useResponsive("up", "lg");

  const isNavMini = lgUp;

  return (
    <Box
      component="main"
      sx={{
        flexGrow: 1,
        minHeight: 1,
        display: "flex",
        flexDirection: "column",
        ...sx,
      }}
      {...other}
    >
      {children}
    </Box>
  );
}
