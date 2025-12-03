#!/usr/bin/env node
/**
 * Pricofy Geocode FR - CDK App Entry Point
 *
 * Defines Lambda function for French postal code geocoding operations.
 */

import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { GeocodeFrStack } from '../lib/geocode-fr-stack';

const app = new cdk.App();

// Get environment from context (defaults to 'dev')
// Support both 'env' and 'environment' for backwards compatibility
const environment = app.node.tryGetContext('env') || app.node.tryGetContext('environment') || 'dev';

// Validate environment
if (!['dev', 'prod'].includes(environment)) {
  throw new Error(`Invalid environment: ${environment}. Must be 'dev' or 'prod'.`);
}

// Common props
const stackProps: cdk.StackProps = {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: 'eu-west-1',
  },
  tags: {
    Project: "Pricofy",
    Service: "Geocode-FR",
    Environment: environment,
  },
};

// Geocode FR Stack: Lambda function for French postal code operations
new GeocodeFrStack(app, `PricofyGeocodeFrStack`, {
  ...stackProps,
  description: `Pricofy Geocode FR (${environment}) - French postal code geocoding Lambda function`,
  environment,
});

console.log(`✅ Stack name: PricofyGeocodeFrStack`);

app.synth();

