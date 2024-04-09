import { Theme } from "@mui/material/styles";

// ----------------------------------------------------------------------

export function card(theme: Theme) {
  return {
    MuiCard: {
      styleOverrides: {
        root: {
          position: "relative",
          // boxShadow: theme.customShadows.card,
          borderRadius: 0,
          zIndex: 0, // Fix Safari overflow: hidden with border radius
          border: "1px solid #00000010",
        },
      },
    },
    MuiCardHeader: {
      styleOverrides: {
        root: {
          padding: theme.spacing(3, 3, 0),
        },
      },
    },
    MuiCardContent: {
      styleOverrides: {
        root: {
          padding: theme.spacing(3),
        },
      },
    },
  };
}
