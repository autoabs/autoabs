/// <reference path="../References.d.ts"/>
import * as React from 'react';
import ConfirmButton from './ConfirmButton';
import * as BuildActions from '../actions/BuildActions';
import * as BuildInfoActions from '../actions/BuildInfoActions';
import * as BuildTypes from '../types/BuildTypes';
import * as MiscUtils from '../utils/MiscUtils';

interface Props {
	build: BuildTypes.Build;
}

interface State {
	locked: boolean;
}

const css = {
	card: {
		flex: '1 0 auto',
		minWidth: '250px',
		maxWidth: '600px',
		height: '113px',
		margin: '5px',
		padding: '0',
	} as React.CSSProperties,
	launch: {
		margin: '10px 10px 0 0',
	} as React.CSSProperties,
	content: {
		padding: '10px 0 10px 10px',
	} as React.CSSProperties,
	name: {
		fontSize: '20px',
	} as React.CSSProperties,
	version: {
		fontSize: '14px',
		margin: '5px 0 0 7px',
	} as React.CSSProperties,
	repo: {
		fontSize: '12px',
		marginTop: '7px',
	} as React.CSSProperties,
	actions: {
		display: 'flex',
		flexDirection: 'row',
		padding: '5px 0',
		justifyContent: 'center',
	} as React.CSSProperties,
	action: {
		width: '80px',
	} as React.CSSProperties,
};

export default class Build extends React.Component<Props, State> {
	constructor(props: Props, context: any) {
		super(props, context);
		this.state = {
			locked: false,
		};
	}

	openDialog = (): void => {
		BuildInfoActions.open(this.props.build.id);
	}

	onArchive = (): void => {
		this.lock();
		BuildActions.archive(this.props.build.id).then(
			this.unlock,
			this.unlock,
		);
	}

	onRebuild = (): void => {
		this.lock();
		BuildActions.rebuild(this.props.build.id).then(
			this.unlock,
			this.unlock,
		);
	}

	lock = (): void => {
		this.setState({
			...this.state,
			locked: true,
		});
	}

	unlock = (): void => {
		this.setState({
			...this.state,
			locked: false,
		});
	}

	render(): JSX.Element {
		let build = this.props.build;

		let start = '-';
		if (build.start !== '0001-01-01T00:00:00Z') {
			start = MiscUtils.formatDate(new Date(build.start));
		}
		let stop = '-';
		if (build.stop !== '0001-01-01T00:00:00Z') {
			stop = MiscUtils.formatDate(new Date(build.stop));
		}

		let actions: JSX.Element[];

		switch (build.state) {
			case 'building':
				actions = [
					<ConfirmButton key="rebuild" label="Rebuild"
						style={css.action}
						className="pt-intent-primary"
						disabled={this.state.locked}
						onConfirm={this.onRebuild}
					/>,
					<ConfirmButton key="archive" label="Archive"
						style={css.action}
						className="pt-intent-danger"
						disabled={this.state.locked}
						onConfirm={this.onArchive}
					/>,
				];
				break;
			case 'pending':
				actions = [
					<ConfirmButton key="archive" label="Archive"
						style={css.action}
						className="pt-intent-danger"
						disabled={this.state.locked}
						onConfirm={this.onArchive}
					/>,
				];
				break;
			case 'failed':
				actions = [
					<ConfirmButton key="rebuild" label="Rebuild"
						style={css.action}
						className="pt-intent-primary"
						disabled={this.state.locked}
						onConfirm={this.onRebuild}
					/>,
					<ConfirmButton key="archive" label="Archive"
						style={css.action}
						className="pt-intent-danger"
						disabled={this.state.locked}
						onConfirm={this.onArchive}
					/>,
				];
				break;
			case 'completed':
				actions = [
					<ConfirmButton key="rebuild" label="Rebuild"
						style={css.action}
						className="pt-intent-primary"
						disabled={this.state.locked}
						onConfirm={this.onRebuild}
					/>,
					<ConfirmButton key="archive" label="Archive"
						style={css.action}
						className="pt-intent-danger"
						disabled={this.state.locked}
						onConfirm={this.onArchive}
					/>,
				];
				break;
			case 'archived':
				actions = [
					<ConfirmButton key="rebuild" label="Rebuild"
						style={css.action}
						className="pt-intent-primary"
						disabled={this.state.locked}
						onConfirm={this.onRebuild}
					/>,
				];
				break;
		}

		return <div className="pt-card" style={css.card}>
			<div className="layout horizontal">
				<div style={css.content} className="card-content flex">
					<div className="layout vertical">
						<div className="layout horizontal">
							<div style={css.name}>{build.name}</div>
							<div  className="pt-text-muted" style={css.version}>
								{build.version}-{build.release} ({build.state})
							</div>
						</div>
						<div className="pt-text-muted" style={css.repo}>
							{build.repo} - {build.arch}
						</div>
					</div>
				</div>
				<div>
					<button type="button"
						className="pt-button pt-minimal pt-icon-document"
						style={css.launch}
						onClick={this.openDialog}
					/>
				</div>
			</div>
			<div className="pt-button-group pt-minimal" style={css.actions}>
				{actions}
			</div>
		</div>;
	}
}
