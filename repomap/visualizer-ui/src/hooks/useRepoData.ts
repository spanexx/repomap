/**
 * Code Map: useRepoData – Custom hook for fetching repomap data and handling file uploads.
 * CID: 1.5.3-useRepoData
 */

import { useState, useEffect, useCallback, useMemo, type ChangeEvent } from 'react';
import { type RepoMapData } from '../types';

export function useRepoData() {
    const [data, setData] = useState<RepoMapData | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Fetch initial data
    useEffect(() => {
        const fetchData = async () => {
            setIsLoading(true);
            try {
                // 1. Try to get the repomap data
                const mapRes = await fetch('/api/repomap');
                if (mapRes.ok) {
                    const mapJson = await mapRes.json();
                    setData(mapJson);
                    setError(null);
                } else {
                    console.warn(`Backend returned ${mapRes.status}`);
                }
            } catch (err) {
                console.warn('Backend not available or data missing', err);
                setError('Failed to fetch data from backend');
            } finally {
                setIsLoading(false);
            }
        };
        fetchData();
    }, []);

    // Handle JSON file upload
    const handleUpload = useCallback((e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        setIsLoading(true);
        const reader = new FileReader();
        reader.onload = (ev) => {
            try {
                const json = JSON.parse(ev.target?.result as string);
                setData(json);
                setError(null);
            } catch (err) {
                console.error("Invalid JSON upload", err);
                alert("Invalid JSON file");
            } finally {
                setIsLoading(false);
            }
        };
        reader.readAsText(file);
    }, []);

    const files = useMemo(() => data?.repomap?.files || data?.files || [], [data]);

    return {
        data,
        files,
        isLoading,
        error,
        handleUpload
    };
}
