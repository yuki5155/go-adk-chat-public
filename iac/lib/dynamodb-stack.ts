import * as cdk from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import { Construct } from 'constructs';

/**
 * DynamoDB Stack Props
 */
export interface DynamoDBStackProps extends cdk.StackProps {
  projectName: string;
  environment: string;
}

/**
 * DynamoDB Stack for User Role Management
 * Creates tables for user roles and role requests
 */
export class DynamoDBStack extends cdk.Stack {
  public readonly userRolesTable: dynamodb.Table;
  public readonly roleRequestsTable: dynamodb.Table;

  constructor(scope: Construct, id: string, props: DynamoDBStackProps) {
    super(scope, id, props);

    const { projectName, environment } = props;

    // ============================================================================
    // User Roles Table
    // Stores user role assignments (excluding root user managed via env var)
    // ============================================================================
    this.userRolesTable = new dynamodb.Table(this, 'UserRolesTable', {
      tableName: `${projectName}-${environment}-user-roles`,
      partitionKey: {
        name: 'user_id',
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST, // On-demand billing for dev
      removalPolicy: environment === 'prod'
        ? cdk.RemovalPolicy.RETAIN
        : cdk.RemovalPolicy.DESTROY,
      pointInTimeRecovery: environment === 'prod', // Enable PITR for production
      encryption: dynamodb.TableEncryption.AWS_MANAGED,
      stream: dynamodb.StreamViewType.NEW_AND_OLD_IMAGES, // For audit logging
    });

    // GSI for email lookup
    this.userRolesTable.addGlobalSecondaryIndex({
      indexName: 'email-index',
      partitionKey: {
        name: 'email',
        type: dynamodb.AttributeType.STRING
      },
      projectionType: dynamodb.ProjectionType.ALL
    });

    // GSI for role lookup
    this.userRolesTable.addGlobalSecondaryIndex({
      indexName: 'role-index',
      partitionKey: {
        name: 'role',
        type: dynamodb.AttributeType.STRING
      },
      projectionType: dynamodb.ProjectionType.ALL
    });

    // ============================================================================
    // Role Requests Table
    // Tracks user requests for role subscriptions
    // ============================================================================
    this.roleRequestsTable = new dynamodb.Table(this, 'RoleRequestsTable', {
      tableName: `${projectName}-${environment}-role-requests`,
      partitionKey: {
        name: 'request_id',
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: environment === 'prod'
        ? cdk.RemovalPolicy.RETAIN
        : cdk.RemovalPolicy.DESTROY,
      pointInTimeRecovery: environment === 'prod',
      encryption: dynamodb.TableEncryption.AWS_MANAGED,
      stream: dynamodb.StreamViewType.NEW_AND_OLD_IMAGES,
    });

    // GSI for user_id lookup
    this.roleRequestsTable.addGlobalSecondaryIndex({
      indexName: 'user_id-index',
      partitionKey: {
        name: 'user_id',
        type: dynamodb.AttributeType.STRING
      },
      projectionType: dynamodb.ProjectionType.ALL
    });

    // GSI for status-based queries with timestamp sorting
    this.roleRequestsTable.addGlobalSecondaryIndex({
      indexName: 'status-requested_at-index',
      partitionKey: {
        name: 'status',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'requested_at',
        type: dynamodb.AttributeType.NUMBER
      },
      projectionType: dynamodb.ProjectionType.ALL
    });

    // ============================================================================
    // CloudFormation Outputs
    // ============================================================================
    new cdk.CfnOutput(this, 'UserRolesTableName', {
      value: this.userRolesTable.tableName,
      description: 'User Roles Table Name',
      exportName: `${projectName}-${environment}-user-roles-table`
    });

    new cdk.CfnOutput(this, 'UserRolesTableArn', {
      value: this.userRolesTable.tableArn,
      description: 'User Roles Table ARN'
    });

    new cdk.CfnOutput(this, 'RoleRequestsTableName', {
      value: this.roleRequestsTable.tableName,
      description: 'Role Requests Table Name',
      exportName: `${projectName}-${environment}-role-requests-table`
    });

    new cdk.CfnOutput(this, 'RoleRequestsTableArn', {
      value: this.roleRequestsTable.tableArn,
      description: 'Role Requests Table ARN'
    });

    // Console output
    console.log('');
    console.log('=========================');
    console.log('✅ DynamoDB Tables Created');
    console.log('=========================');
    console.log(`📊 User Roles Table: ${this.userRolesTable.tableName}`);
    console.log('   - PK: user_id');
    console.log('   - GSI: email-index');
    console.log('   - GSI: role-index');
    console.log('');
    console.log(`📊 Role Requests Table: ${this.roleRequestsTable.tableName}`);
    console.log('   - PK: request_id');
    console.log('   - GSI: user_id-index');
    console.log('   - GSI: status-requested_at-index');
    console.log('=========================');
  }
}
