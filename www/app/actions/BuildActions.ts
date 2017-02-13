/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import * as Alert from '../Alert';
import Loader from '../Loader';
import * as BuildTypes from '../types/BuildTypes';
import BuildStore from '../stores/BuildStore';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;

function _sync(index: number): Promise<string> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<string>((resolve, reject): void => {
		SuperAgent
			.get('/build')
			.query({'index': index})
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (curSyncId !== syncId) {
					return;
				}

				if (err) {
					Alert.error('Failed to sync builds');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: BuildTypes.SYNC,
					data: {
						builds: res.body.builds,
						index: res.body.index,
						count: res.body.count,
					},
				});

				resolve();
			});
	});
}

export function traverse(index: number): Promise<string> {
	BuildStore._index = index;
	return _sync(index);
}

export function sync(): Promise<string> {
	return _sync(BuildStore.index);
}

export function archive(id: string): Promise<string> {
	let loader = new Loader().loading();

	return new Promise<string>((resolve, reject): void => {
		SuperAgent
			.put('/build/' + id + '/archive')
			.set('Accept', 'application/json')
			.end((err: any): void => {
				loader.done();

				if (err) {
					Alert.error('Failed to archive build');
					reject(err);
					return;
				}

				sync().then(resolve, resolve);
			});
	});
}

export function rebuild(id: string): Promise<string> {
	let loader = new Loader().loading();

	return new Promise<string>((resolve, reject): void => {
		SuperAgent
			.put('/build/' + id + '/rebuild')
			.set('Accept', 'application/json')
			.end((err: any): void => {
				loader.done();

				if (err) {
					Alert.error('Failed to rebuild build');
					reject(err);
					return;
				}

				sync().then(resolve, resolve);
			});
	});
}

export function remove(id: string): void {
	Dispatcher.dispatch({
		type: BuildTypes.REMOVE,
		data: {
			id: id,
		},
	});

	Alert.info('Build successfully removed');
}
