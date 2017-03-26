/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import * as Events from 'events';
import * as NodeTypes from '../types/NodeTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class NodeStore extends Events.EventEmitter {
	_nodes: NodeTypes.Nodes = [];
	_token = Dispatcher.register((this._callback).bind(this));

	get nodes(): NodeTypes.Nodes {
		return this._nodes;
	}

	emitChange(): void {
		this.emit(GlobalTypes.CHANGE);
	}

	addChangeListener(callback: () => void): void {
		this.on(GlobalTypes.CHANGE, callback);
	}

	removeChangeListener(callback: () => void): void {
		this.removeListener(GlobalTypes.CHANGE, callback);
	}

	_sync(nodes: NodeTypes.Nodes): void {
		this._nodes = nodes;
		this.emitChange();
	}

	_callback(action: NodeTypes.NodeDispatch): void {
		switch (action.type) {
			case NodeTypes.SYNC:
				this._sync(action.data.nodes);
				break;
		}
	}
}

export default new NodeStore();
