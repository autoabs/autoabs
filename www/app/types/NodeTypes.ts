/// <reference path="../References.d.ts"/>
export const SYNC = 'node.sync';
export const CHANGE = 'node.change';

export interface Node {
	id: string;
	type: string;
	memory: number;
	load1: number;
	load5: number;
	load15: number;
	stats: NodeStats;
	settings: NodeSettings;
}

export type NodeStats = DefaultStats | BuilderStats;

export type DefaultStats = {}

export interface BuilderStats {
	active: number;
}

export type NodeSettings = DefaultSettings | BuilderSettings;

export type DefaultSettings = {}

export interface BuilderSettings {
	concurrency: number;
}

export type Nodes = Node[];

export interface NodeDispatch {
	type: string;
	data?: {
		nodes?: Node[];
	};
}
