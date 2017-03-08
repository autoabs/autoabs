/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import Loader from '../Loader';
import * as BuildTypes from '../types/BuildTypes';
import BuildStore from '../stores/BuildStore';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;

function _sync(index: number): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/build')
			.query({'index': index})
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (curSyncId !== syncId) {
					resolve();
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

export function traverse(index: number): Promise<void> {
	BuildStore._index = index;
	return _sync(index);
}

export function sync(): Promise<void> {
	return _sync(BuildStore.index);
}

export function archive(id: string): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
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

export function rebuild(id: string): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
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

EventDispatcher.register((action: BuildTypes.BuildDispatch) => {
	switch (action.type) {
		case BuildTypes.CHANGE:
			sync();
			break;
	}
});
