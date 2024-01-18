import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.tsx';
import '@mantine/core/styles.css';
import '@mantine/dates/styles.css';
import '@mantine/charts/styles.css';
import 'dayjs/locale/th';
import { MantineProvider } from '@mantine/core';
import { DatesProvider } from '@mantine/dates';
import { Provider as JotaiProvider } from 'jotai';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <JotaiProvider>
      <MantineProvider>
        <DatesProvider settings={{ locale: 'th' }}>
          <App />
        </DatesProvider>
      </MantineProvider>
    </JotaiProvider>
  </React.StrictMode>
);
