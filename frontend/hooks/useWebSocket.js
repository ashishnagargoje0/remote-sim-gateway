import { useState, useEffect, useRef, useContext, createContext } from 'react';
import { useAuth } from './useAuth';
import toast from 'react-hot-toast';

const WebSocketContext = createContext({});

export function WebSocketProvider({ children }) {
  const [connectionStatus, setConnectionStatus] = useState('disconnected');
  const [connectedDevices, setConnectedDevices] = useState([]);
  const [messages, setMessages] = useState([]);
  const socketRef = useRef(null);
  const reconnectTimeoutRef = useRef(null);
  const { user, isAuthenticated } = useAuth();

  const connect = () => {
    if (!isAuthenticated || !user) return;

    try {
      const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws';
      const socket = new WebSocket(`${wsUrl}?user_id=${user.id}`);

      socket.onopen = () => {
        console.log('WebSocket connected');
        setConnectionStatus('connected');
        toast.success('Connected to gateway');
        
        // Clear any existing reconnect timeout
        if (reconnectTimeoutRef.current) {
          clearTimeout(reconnectTimeoutRef.current);
          reconnectTimeoutRef.current = null;
        }
      };

      socket.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          handleMessage(message);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      socket.onclose = () => {
        console.log('WebSocket disconnected');
        setConnectionStatus('disconnected');
        
        // Attempt to reconnect after 5 seconds
        if (isAuthenticated) {
          reconnectTimeoutRef.current = setTimeout(() => {
            console.log('Attempting to reconnect...');
            connect();
          }, 5000);
        }
      };

      socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        setConnectionStatus('error');
        toast.error('Connection error');
      };

      socketRef.current = socket;
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
      setConnectionStatus('error');
    }
  };

  const disconnect = () => {
    if (socketRef.current) {
      socketRef.current.close();
      socketRef.current = null;
    }
    
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
    
    setConnectionStatus('disconnected');
    setConnectedDevices([]);
  };

  const sendMessage = (message) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify(message));
      return true;
    }
    return false;
  };

  const handleMessage = (message) => {
    console.log('Received WebSocket message:', message);

    switch (message.type) {
      case 'welcome':
        console.log('Welcome message received:', message.data);
        break;

      case 'device_connected':
        setConnectedDevices(prev => {
          const deviceId = message.data.device_id;
          if (!prev.find(d => d.device_id === deviceId)) {
            toast.success(`Device ${message.data.name || deviceId} connected`);
            return [...prev, message.data];
          }
          return prev;
        });
        break;

      case 'device_disconnected':
        setConnectedDevices(prev => {
          const filtered = prev.filter(d => d.device_id !== message.data.device_id);
          if (filtered.length !== prev.length) {
            toast.error(`Device ${message.data.name || message.data.device_id} disconnected`);
          }
          return filtered;
        });
        break;

      case 'sms_status_update':
        toast.success(`SMS ${message.data.status}: ${message.data.phone_number}`);
        break;

      case 'call_status_update':
        toast.info(`Call ${message.data.status}: ${message.data.phone_number}`);
        break;

      case 'device_status_update':
        setConnectedDevices(prev => 
          prev.map(device => 
            device.device_id === message.data.device_id 
              ? { ...device, ...message.data }
              : device
          )
        );
        break;

      case 'error':
        toast.error(message.data.message || 'An error occurred');
        break;

      default:
        console.log('Unknown message type:', message.type);
    }

    // Store message in history
    setMessages(prev => [message, ...prev.slice(0, 99)]); // Keep last 100 messages
  };

  // Connect when authenticated
  useEffect(() => {
    if (isAuthenticated && user) {
      connect();
    } else {
      disconnect();
    }

    // Cleanup on unmount
    return () => {
      disconnect();
    };
  }, [isAuthenticated, user]);

  // Cleanup timeouts on unmount
  useEffect(() => {
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
    };
  }, []);

  const value = {
    connectionStatus,
    connectedDevices,
    messages,
    sendMessage,
    connect,
    disconnect,
  };

  return (
    <WebSocketContext.Provider value={value}>
      {children}
    </WebSocketContext.Provider>
  );
}

export function useWebSocket() {
  const context = useContext(WebSocketContext);
  if (context === undefined) {
    throw new Error('useWebSocket must be used within a WebSocketProvider');
  }
  return context;
}