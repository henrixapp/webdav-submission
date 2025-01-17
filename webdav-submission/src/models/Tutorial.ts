// tslint:disable
/**
 * Submissions
 * This API specifies the submissions service, as accessed by the web admin UI used by students, lecturers and tutors. 
 *
 * The version of the OpenAPI document: 1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
/**
 * The tutorial is part of an lecture and can have multiple tutors
 * @export
 * @interface Tutorial
 */
export interface Tutorial  {
    /**
     * The globally recognized ID of an Tutorial.
     * @type {string}
     * @memberof Tutorial
     */
    id?: string;
    /**
     * The ID of the lecture this tutorial is located in. (foreign)
     * @type {number}
     * @memberof Tutorial
     */
    lectureID?: number;
    /**
     * A short descriptive title of the tutorial
     * @type {string}
     * @memberof Tutorial
     */
    title?: string;
    /**
     * 
     * @type {Date}
     * @memberof Tutorial
     */
    createdAt?: Date;
    /**
     * 
     * @type {Date}
     * @memberof Tutorial
     */
    updatedAt?: Date;
}

export function TutorialFromJSON(json: any): Tutorial {
    return {
        'id': !exists(json, 'id') ? undefined : json['id'],
        'lectureID': !exists(json, 'lectureID') ? undefined : json['lectureID'],
        'title': !exists(json, 'title') ? undefined : json['title'],
        'createdAt': !exists(json, 'createdAt') ? undefined : new Date(json['createdAt']),
        'updatedAt': !exists(json, 'updatedAt') ? undefined : new Date(json['updatedAt']),
    };
}

export function TutorialToJSON(value?: Tutorial): any {
    if (value === undefined) {
        return undefined;
    }
    return {
        'id': value.id,
        'lectureID': value.lectureID,
        'title': value.title,
        'createdAt': value.createdAt === undefined ? undefined : value.createdAt.toISOString(),
        'updatedAt': value.updatedAt === undefined ? undefined : value.updatedAt.toISOString(),
    };
}


