import { FieldGroup } from 'models/Field';
import { MockTag } from 'models/Tag';
import { TagResource } from 'models/TagResource';

export interface TagState {
	tagResource: TagResource;
	tags: MockTag[];
	fields: FieldGroup[];
}

export interface ErrorType {
	response?: { status: string };
	message: string;
}
