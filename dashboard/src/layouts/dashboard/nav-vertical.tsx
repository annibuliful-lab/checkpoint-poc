import { useCallback, useEffect } from "react";

import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";
import Drawer from "@mui/material/Drawer";

import { useNavData } from "./config-navigation";
import { useResponsive } from "@/hooks/use-responsive";
import { usePathname } from "next/navigation";
import { NAV } from "../config-layout";
import { NavSectionVertical } from "@/components/nav-section";
import Scrollbar from "@/components/scrollbar";
import Logo from "@/components/logo";
import { IconButton } from "@mui/material";
import { Logout } from "@mui/icons-material";
import { useSetAtom } from "jotai";
import { authAtom } from "@/auth/store";
import { apolloCache } from "@/apollo-client/apollo-wrapper";

// ----------------------------------------------------------------------

type Props = {
  openNav: boolean;
  onCloseNav: VoidFunction;
};

export default function NavVertical({ openNav, onCloseNav }: Props) {
  const user = {
    role: "USER",
  };

  const pathname = usePathname();

  const lgUp = useResponsive("up", "lg");

  const navData = useNavData();
  const setAuth = useSetAtom(authAtom);
  const handleLogout = useCallback(async () => {
    setAuth(undefined);
    await apolloCache.reset();
    window.localStorage.clear();
  }, [setAuth]);

  useEffect(() => {
    if (openNav) {
      onCloseNav();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pathname]);

  const renderContent = (
    <Scrollbar
      sx={{
        height: 1,
        "& .simplebar-content": {
          height: 1,
          display: "flex",
          flexDirection: "column",
        },
      }}
    >
      <Logo sx={{ mt: 3, ml: 3, mb: 1 }} />

      <NavSectionVertical
        data={navData}
        slotProps={{
          currentRole: user?.role,
        }}
      />

      <Box sx={{ flexGrow: 1 }} />
      <Box sx={{ mx: "auto", mb: 1 }}>
        <IconButton onClick={handleLogout}>
          <Logout />
        </IconButton>
      </Box>
    </Scrollbar>
  );

  return (
    <Box
      sx={{
        flexShrink: { lg: 0 },
        width: { lg: NAV.W_MINI },
      }}
    >
      {lgUp ? (
        <Stack
          sx={{
            height: 1,
            position: "fixed",
            width: NAV.W_MINI,
            borderRight: (theme) => `solid 1px ${theme.palette.divider}`,
          }}
        >
          {renderContent}
        </Stack>
      ) : (
        <Drawer
          open={openNav}
          onClose={onCloseNav}
          PaperProps={{
            sx: {
              width: NAV.W_MINI,
            },
          }}
        >
          {renderContent}
        </Drawer>
      )}
    </Box>
  );
}
