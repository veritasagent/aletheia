"use client";

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
  type PropsWithChildren,
} from "react";

export interface ConnectionState {
  connected: boolean;
  total: number;
  falseCount: number;
  chainCount: number;
}

interface ConnectionContextValue extends ConnectionState {
  setConnected: (connected: boolean) => void;
  setMetrics: (metrics: Partial<Omit<ConnectionState, "connected">>) => void;
  setConnectionState: (state: ConnectionState) => void;
  resetConnectionState: () => void;
}

const initialState: ConnectionState = {
  connected: false,
  total: 0,
  falseCount: 0,
  chainCount: 0,
};

export const ConnectionContext = createContext<ConnectionContextValue | undefined>(undefined);

export function ConnectionProvider({ children }: PropsWithChildren) {
  const [state, setState] = useState<ConnectionState>(initialState);

  const setConnected = useCallback((connected: boolean) => {
    setState((prev) => ({ ...prev, connected }));
  }, []);

  const setMetrics = useCallback((metrics: Partial<Omit<ConnectionState, "connected">>) => {
    setState((prev) => ({
      ...prev,
      total: metrics.total ?? prev.total,
      falseCount: metrics.falseCount ?? prev.falseCount,
      chainCount: metrics.chainCount ?? prev.chainCount,
    }));
  }, []);

  const setConnectionState = useCallback((nextState: ConnectionState) => {
    setState(nextState);
  }, []);

  const resetConnectionState = useCallback(() => {
    setState(initialState);
  }, []);

  const value = useMemo<ConnectionContextValue>(
    () => ({
      ...state,
      setConnected,
      setMetrics,
      setConnectionState,
      resetConnectionState,
    }),
    [state, setConnected, setMetrics, setConnectionState, resetConnectionState],
  );

  return <ConnectionContext.Provider value={value}>{children}</ConnectionContext.Provider>;
}

export function useConnectionContext(): ConnectionContextValue {
  const context = useContext(ConnectionContext);
  if (!context) {
    throw new Error("useConnectionContext must be used inside ConnectionProvider");
  }
  return context;
}
