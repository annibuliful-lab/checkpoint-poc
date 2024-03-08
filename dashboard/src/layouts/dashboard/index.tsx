import Box from "@mui/material/Box";

import Main from "./main";
import NavVertical from "./nav-vertical";
import { useBoolean } from "@/hooks/use-boolean";

// ----------------------------------------------------------------------

type Props = {
  children: React.ReactNode;
};

export default function DashboardLayout({ children }: Props) {
  const nav = useBoolean();

  const renderNavVertical = (
    <NavVertical openNav={nav.value} onCloseNav={nav.onFalse} />
  );

  return (
    <>
      <Box
        sx={{
          minHeight: 1,
          display: "flex",
          flexDirection: { xs: "column", lg: "row" },
        }}
      >
        {renderNavVertical}

        <Main>{children}</Main>
      </Box>
    </>
  );
}
